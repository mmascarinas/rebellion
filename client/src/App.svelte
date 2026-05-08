<script>
  import Router from './lib/Router.svelte'
  import Home from './pages/Home/Home.svelte'
  import Lobby from './pages/Lobby/Lobby.svelte'
  import Game from './pages/Game/Game.svelte'

  const routes = {
    '/': { component: Home, showHeader: false },
    '/lobby': { component: Lobby },
    '/game': { component: Game },
  }
</script>

<Router {routes} />
