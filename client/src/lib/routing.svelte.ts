// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface RouteConfig {
	component: any
	showHeader?: boolean
}

export type RouteMap = Record<string, RouteConfig>

let current = $state(window.location.hash.slice(1) || '/')

window.addEventListener('hashchange', () => {
	current = window.location.hash.slice(1) || '/'
})

export function currentRoute() {
	return current
}

export function push(path: string) {
	window.location.hash = path
}
