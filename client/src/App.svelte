<script>
	import Router from './lib/Router.svelte'
	import Home from './pages/Home/Home.svelte'
	import Lobby from './pages/Lobby/Lobby.svelte'
	import Game from './pages/Game/Game.svelte'
	import { name } from '../package.json'

	const title = name.charAt(0).toUpperCase() + name.slice(1)

	const routes = {
		'/': { component: Home, showHeader: false, title: `${title} - Home` },
		'/lobby': { component: Lobby, title: `${title} - Lobby` },
		'/game': { component: Game, title: `${title} - Game` }
	}
</script>

<Router {routes} />
