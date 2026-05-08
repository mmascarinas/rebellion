export function setTheme(theme: 'light' | 'dark') {
	document.documentElement.setAttribute('data-theme', theme)
}

export function toggleTheme() {
	const current = document.documentElement.getAttribute('data-theme')
	setTheme(current === 'dark' ? 'light' : 'dark')
}
