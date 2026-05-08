let current = $state(
	(document.documentElement.getAttribute('data-theme') as 'light' | 'dark') ??
		'light'
)

export function isDark() {
	return current === 'dark'
}

export function setTheme(theme: 'light' | 'dark') {
	document.documentElement.setAttribute('data-theme', theme)
	current = theme
}

export function toggleTheme() {
	setTheme(current === 'dark' ? 'light' : 'dark')
}
