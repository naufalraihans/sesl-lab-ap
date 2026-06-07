import { test, expect } from '@playwright/test';

test('login page has expected title and can be interacted with', async ({ page }) => {
	await page.goto('/praktikum/login');

	// Expect a title "to contain" a substring.
	await expect(page).toHaveTitle(/Lab Algoritma/);

	// Expect the NIM input to be visible
	const nimInput = page.getByPlaceholder('Masukkan NIM Anda');
	await expect(nimInput).toBeVisible();

	// We can type into it
	await nimInput.fill('1202220000');
	await expect(nimInput).toHaveValue('1202220000');
	
	// Expect the button to exist
	const loginBtn = page.getByRole('button', { name: /Lanjutkan/i });
	await expect(loginBtn).toBeVisible();
});
