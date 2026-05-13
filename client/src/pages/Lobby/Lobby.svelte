<script lang="ts">
	import { push } from '../../lib/routing.svelte'
	import CreateGameModal from '../../components/CreateGameModal.svelte'
	import SearchBar from '../../components/SearchBar.svelte'
	import GameCard, { type Game } from '../../components/GameCard.svelte'
	import FilterGroup from '../../components/FilterGroup.svelte'
	import FilterSelect from '../../components/FilterSelect.svelte'
	import StatsBar, { type Stat } from '../../components/StatsBar.svelte'
	import Pagination from '../../components/Pagination.svelte'

	let showCreateGame = $state(false)

	type GameStatus = Game['status']
	type GameVisibility = Game['visibility']

	interface GameListResponse {
		total: number
		perPage: number
		page: number
		results: Game[]
	}

	const mockGames: GameListResponse = {
		total: 100,
		perPage: 10,
		page: 1,
		results: [
			{
				id: '1',
				name: 'Royal Court',
				host: 'Emperor_X',
				players: 3,
				maxPlayers: 6,
				status: 'waiting',
				visibility: 'public',
				createdAgo: '2m ago'
			},
			{
				id: '2',
				name: 'Assassination Room',
				host: 'ShadowBlade',
				players: 5,
				maxPlayers: 5,
				status: 'active',
				visibility: 'public',
				createdAgo: '8m ago'
			},
			{
				id: '3',
				name: 'The Vault',
				host: 'CryptoKing',
				players: 2,
				maxPlayers: 4,
				status: 'waiting',
				visibility: 'private',
				createdAgo: '1m ago'
			},
			{
				id: '4',
				name: 'Bluff Masters',
				host: 'PokerFace99',
				players: 4,
				maxPlayers: 4,
				status: 'full',
				visibility: 'public',
				createdAgo: '12m ago'
			},
			{
				id: '5',
				name: 'Coup de Grâce',
				host: 'Strategist',
				players: 1,
				maxPlayers: 6,
				status: 'waiting',
				visibility: 'public',
				createdAgo: 'just now'
			},
			{
				id: '6',
				name: 'Inner Circle',
				host: 'VIP_Only',
				players: 3,
				maxPlayers: 4,
				status: 'active',
				visibility: 'private',
				createdAgo: '5m ago'
			},
			{
				id: '7',
				name: 'No Mercy',
				host: 'Ruthless',
				players: 2,
				maxPlayers: 6,
				status: 'waiting',
				visibility: 'public',
				createdAgo: '30s ago'
			},
			{
				id: '8',
				name: 'Beginners Welcome',
				host: 'FriendlyHost',
				players: 4,
				maxPlayers: 6,
				status: 'waiting',
				visibility: 'public',
				createdAgo: '4m ago'
			},
			{
				id: '9',
				name: 'Elite Arena',
				host: 'ProGamer',
				players: 6,
				maxPlayers: 6,
				status: 'active',
				visibility: 'private',
				createdAgo: '15m ago'
			}
		]
	}

	let searchQuery = $state('')
	let filterStatus = $state<GameStatus | 'all'>('all')
	let filterVisibility = $state<GameVisibility | 'all'>('all')
	let currentPage = $state(1)

	const totalPages = $derived(Math.ceil(mockGames.total / mockGames.perPage))

	const filteredGames = $derived(
		mockGames.results.filter((game) => {
			const matchesSearch =
				game.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				game.host.toLowerCase().includes(searchQuery.toLowerCase())
			const matchesStatus =
				filterStatus === 'all' || game.status === filterStatus
			const matchesVisibility =
				filterVisibility === 'all' || game.visibility === filterVisibility
			return matchesSearch && matchesStatus && matchesVisibility
		})
	)

	const stats = $derived({
		total: mockGames.total,
		waiting: mockGames.results.filter((g) => g.status === 'waiting').length,
		active: mockGames.results.filter((g) => g.status === 'active').length,
		players: mockGames.results.reduce((sum, g) => sum + g.players, 0)
	})
</script>

<div class="lobby">
	<div class="lobby-header">
		<div class="lobby-title-row">
			<h1 class="lobby-title">Game Lobby</h1>
			<button class="create-btn" onclick={() => (showCreateGame = true)}>
				<span class="create-icon">+</span>
				Create Game
			</button>
		</div>

		<StatsBar
			stats={[
				{ value: stats.total, label: 'Games' },
				{ value: stats.waiting, label: 'Waiting' },
				{ value: stats.active, label: 'In Progress' },
				{ value: stats.players, label: 'Players Online' }
			]}
		/>
	</div>

	<div class="toolbar">
		<SearchBar
			bind:value={searchQuery}
			placeholder="Search games or hosts..."
		/>

		<div class="filters">
			<FilterGroup
				bind:value={filterStatus}
				options={[
					{ value: 'all', label: 'All' },
					{ value: 'waiting', label: 'Waiting' },
					{ value: 'active', label: 'Active' },
					{ value: 'full', label: 'Full' }
				]}
			/>

			<FilterGroup
				bind:value={filterVisibility}
				options={[
					{ value: 'all', label: 'All' },
					{ value: 'public', label: 'Public' },
					{ value: 'private', label: 'Private' }
				]}
			/>
		</div>

		<div class="filters-mobile">
			<FilterSelect
				bind:value={filterStatus}
				options={[
					{ value: 'all', label: 'All Status' },
					{ value: 'waiting', label: 'Waiting' },
					{ value: 'active', label: 'Active' },
					{ value: 'full', label: 'Full' }
				]}
			/>

			<FilterSelect
				bind:value={filterVisibility}
				options={[
					{ value: 'all', label: 'All Visibility' },
					{ value: 'public', label: 'Public' },
					{ value: 'private', label: 'Private' }
				]}
			/>
		</div>
	</div>

	<div class="game-grid">
		{#each filteredGames as game (game.id)}
			<GameCard {game} />
		{:else}
			<div class="empty-state">
				<p class="empty-icon">⌕</p>
				<p class="empty-text">No games found</p>
				<p class="empty-sub">Try adjusting your filters</p>
			</div>
		{/each}
	</div>

	<Pagination bind:currentPage {totalPages} />
</div>

{#if showCreateGame}
	<CreateGameModal
		onclose={() => (showCreateGame = false)}
		oncreate={(config) => {
			console.log('Create game:', config)
			showCreateGame = false
		}}
	/>
{/if}

<style lang="scss">
	.lobby {
		display: flex;
		flex-direction: column;
		gap: 24px;
		padding: 32px clamp(24px, 5vw, 64px);
		max-width: 1400px;
		margin: 0 auto;
		width: 100%;
	}

	/* Header */
	.lobby-header {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.lobby-title-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.lobby-title {
		font-size: clamp(20px, 3vw, 28px);
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-text);
	}

	.create-btn {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 10px 20px;
		background: var(--color-primary);
		color: #fff;
		border: none;
		border-radius: 8px;
		font-family: inherit;
		font-size: 13px;
		font-weight: 500;
		letter-spacing: 0.05em;
		cursor: pointer;
		transition:
			transform 0.15s ease,
			box-shadow 0.15s ease;

		&:hover {
			transform: translateY(-1px);
			box-shadow: 0 4px 16px
				color-mix(in srgb, var(--color-primary) 40%, transparent);
		}
	}

	.create-icon {
		font-size: 18px;
		font-weight: 300;
		line-height: 1;
	}

	/* Toolbar */
	.toolbar {
		display: flex;
		align-items: center;
		gap: 16px;
		flex-wrap: wrap;
	}

	.filters {
		display: flex;
		gap: 8px;
	}

	.filters-mobile {
		display: none;
		gap: 8px;
		width: 100%;
	}

	/* Game Grid */
	.game-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
		gap: 16px;
	}

	/* Empty state */
	.empty-state {
		grid-column: 1 / -1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		padding: 64px 24px;
	}

	.empty-icon {
		font-size: 48px;
		color: var(--color-text-secondary);
		opacity: 0.3;
	}

	.empty-text {
		font-size: 16px;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.empty-sub {
		font-size: 13px;
		color: var(--color-text-secondary);
		opacity: 0.6;
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.4;
		}
	}

	@media (max-width: 550px) {
		.lobby {
			padding: 20px 16px;
			gap: 16px;
		}

		.lobby-title-row {
			flex-direction: column;
			align-items: stretch;
			gap: 12px;
		}

		.toolbar {
			flex-direction: column;
			align-items: stretch;
		}

		.filters {
			width: 100%;
			overflow-x: auto;
		}

		.game-grid {
			grid-template-columns: 1fr;
		}

		.create-btn {
			justify-content: center;
		}
	}

	@media (max-width: 460px) {
		.filters {
			display: none;
		}

		.filters-mobile {
			display: flex;
		}
	}
</style>
