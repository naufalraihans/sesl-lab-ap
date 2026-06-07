import type { PageLoad } from './$types';

export const load: PageLoad = async ({ setHeaders }) => {
	// Cache the public lobby page for 5 minutes (300 seconds) in the browser
	setHeaders({
		'cache-control': 'public, max-age=300'
	});
};
