import os
import sys
from PIL import Image

# Coba import pillow-heif untuk mendukung konversi foto iPhone (.HEIC)
try:
    from pillow_heif import register_heif_opener
    register_heif_opener()
    HEIC_SUPPORT = True
except ImportError:
    HEIC_SUPPORT = False

def compress_images(input_dir, output_dir, max_size=1920, quality=85):
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    print("=" * 60)
    print("                Lab AP - Image Compressor for Web              ")
    print("=" * 60)
    if HEIC_SUPPORT:
        print("[✓] Dukungan HEIC/iPhone Aktif!")
    else:
        print("[!] Info: Instal 'pillow-heif' untuk konversi foto HEIC/iPhone.")
        print("    Cara instal: pip install pillow-heif")
    print(f"Target Direktori : {input_dir}")
    print(f"Output Direktori : {output_dir}")
    print(f"Lebar Maksimal   : {max_size}px")
    print(f"Kualitas Gambar  : {quality}%")
    print("-" * 60)

    # Format file yang didukung
    valid_extensions = ('.png', '.jpg', '.jpeg', '.webp', '.gif')
    if HEIC_SUPPORT:
        valid_extensions += ('.heic', '.HEIC')
    else:
        heic_files = [f for f in os.listdir(input_dir) if f.lower().endswith('.heic')]
        if heic_files:
            print(f"[!] Ditemukan {len(heic_files)} file HEIC, tetapi dilewati karena 'pillow-heif' belum diinstal.")

    files = [f for f in os.listdir(input_dir) if f.lower().endswith(valid_extensions)]
    
    if not files:
        print("Tidak ada foto yang ditemukan untuk dikompres.")
        return

    success_count = 0
    total_original_size = 0
    total_compressed_size = 0
    files_to_delete = []

    for idx, filename in enumerate(files, 1):
        input_path = os.path.join(input_dir, filename)
        
        # Selalu simpan ke format .jpg terkompresi
        name_without_ext = os.path.splitext(filename)[0]
        output_filename = f"{name_without_ext}.jpg"
        output_path = os.path.join(output_dir, output_filename)
        
        try:
            original_size = os.path.getsize(input_path)
            total_original_size += original_size

            # Buka gambar
            img = Image.open(input_path)
            
            # Perbaiki orientasi EXIF (misal foto tegak dari kamera HP agar tidak miring)
            try:
                from PIL import ImageOps
                img = ImageOps.exif_transpose(img)
            except Exception:
                pass

            # Hitung dimensi baru (jaga rasio aspek)
            width, height = img.size
            if width > max_size:
                ratio = max_size / float(width)
                new_height = int(float(height) * float(ratio))
                img = img.resize((max_size, new_height), Image.Resampling.LANCZOS)
                print(f"[{idx}/{len(files)}] Resizing {filename} ({width}x{height} -> {max_size}x{new_height})...")
            else:
                print(f"[{idx}/{len(files)}] Kompresi {filename} (resolusi sudah aman: {width}x{height})...")

            # Konversi transparansi PNG/GIF ke warna latar belakang putih
            if img.mode in ('RGBA', 'LA') or (img.mode == 'P' and 'transparency' in img.info):
                bg = Image.new('RGB', img.size, (255, 255, 255))
                bg.paste(img, mask=img.convert('RGBA').split()[3])
                img = bg
            elif img.mode != 'RGB':
                img = img.convert('RGB')

            # Simpan gambar sebagai JPEG terkompresi
            img.save(output_path, 'JPEG', quality=quality, optimize=True)
            
            compressed_size = os.path.getsize(output_path)
            total_compressed_size += compressed_size
            
            saving = (original_size - compressed_size) / original_size * 100
            print(f"    -> Sukses! Ukuran: {original_size/1024/1024:.2f}MB -> {compressed_size/1024:.1f}KB (Hemat {saving:.1f}%)")
            success_count += 1
            
            # Jika file aslinya berformat HEIC atau PNG (dan di-overwrite di direktori yang sama), tandai file lama untuk dihapus agar tidak duplikat
            if input_dir == output_dir and filename.lower() != output_filename.lower():
                files_to_delete.append(input_path)
            
        except Exception as e:
            print(f"[X] Gagal memproses {filename}: {e}")

    # Hapus file lama yang tidak terpakai jika dikompresi di folder yang sama (misal HEIC yang sudah jadi JPG)
    if files_to_delete:
        print("\n[~] Membersihkan file duplikat format lama (.heic/.png)...")
        for path in files_to_delete:
            try:
                os.remove(path)
                print(f"    -> Dihapus: {os.path.basename(path)}")
            except Exception as e:
                print(f"    -> Gagal menghapus {os.path.basename(path)}: {e}")

    print("-" * 60)
    print("                       RINGKASAN KOMPRESI                       ")
    print("-" * 60)
    print(f"Total Foto Berhasil : {success_count} / {len(files)}")
    if success_count > 0:
        print(f"Ukuran Sebelum      : {total_original_size/1024/1024:.2f} MB")
        print(f"Ukuran Sesudah       : {total_compressed_size/1024/1024:.2f} MB")
        saving_total = (total_original_size - total_compressed_size) / total_original_size * 100
        print(f"Total Ruang Dihemat : {saving_total:.1f}% space di disk!")
    print("=" * 60)

if __name__ == '__main__':
    # Secara default menargetkan folder galeri proyek
    default_input = 'frontend/src/lib/galeri'
    default_output = 'frontend/src/lib/galeri'
    
    # Jika dijalankan via terminal dengan argumen custom
    if len(sys.argv) > 1:
        default_input = sys.argv[1]
    if len(sys.argv) > 2:
        default_output = sys.argv[2]
        
    compress_images(default_input, default_output)
