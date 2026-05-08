<script lang="ts">
  import { currentRoute, type RouteMap } from './routing.svelte'
  import Header from '../components/Header.svelte'

  let { routes }: { routes: RouteMap } = $props()

  const config = $derived(routes[currentRoute()] ?? null)
  const Page = $derived(config?.component ?? null)
  const showHeader = $derived(config?.showHeader ?? true)
</script>

{#if Page}
  {#if showHeader}
    <Header />
  {/if}
  <Page />
{:else}
  <h1>404</h1>
  <p>Page not found.</p>
{/if}
