// Format sisa waktu (detik) menjadi mm:ss.
export function fmtCountdown(totalSeconds: number): string {
	if (totalSeconds < 0) totalSeconds = 0;
	const m = Math.floor(totalSeconds / 60);
	const s = totalSeconds % 60;
	return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
}

// Hitung sisa detik sampai deadline ISO string.
export function secondsUntil(deadlineISO?: string | null): number {
	if (!deadlineISO) return 0;
	const diff = new Date(deadlineISO).getTime() - Date.now();
	return Math.max(0, Math.floor(diff / 1000));
}

export function labelJenis(jenis: string): string {
	switch (jenis) {
		case 'pretest': return 'Pre-test';
		case 'posttest': return 'Post-test';
		case 'keterampilan': return 'Keterampilan';
		case 'ujian_praktik': return 'Ujian Praktik';
		default: return jenis;
	}
}

export function labelStatus(status: string): string {
	switch (status) {
		case 'belum_dikerjakan': return 'Belum Dikerjakan';
		case 'sedang_dikerjakan': return 'Sedang Dikerjakan';
		case 'selesai': return 'Selesai';
		default: return status;
	}
}

export function statusBadgeClass(status: string): string {
	switch (status) {
		case 'selesai': return 'bg-state-success-bg text-state-success';
		case 'sedang_dikerjakan': return 'bg-state-warning-bg text-state-warning';
		default: return 'bg-gray-100 text-ink-caption';
	}
}
