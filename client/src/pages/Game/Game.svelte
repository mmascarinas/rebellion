<script lang="ts">
	import TableCenter from '../../components/TableCenter.svelte'
	import PlayerSeat from '../../components/PlayerSeat.svelte'
	import ActionBar from '../../components/ActionBar.svelte'
	import GameLog from '../../components/GameLog.svelte'

	type CardRole = 'courtesan' | 'mercenary' | 'baron' | 'consul' | 'admiral'

	interface PlayerCard {
		role: CardRole
		revealed: boolean
	}

	interface Player {
		id: string
		name: string
		avatar: string
		coins: number
		cards: PlayerCard[]
		isCurrentTurn: boolean
		isEliminated: boolean
		isYou: boolean
	}

	const players: Player[] = [
		{
			id: '1',
			name: 'You',
			avatar: 'Y',
			coins: 4,
			cards: [
				{ role: 'baron', revealed: false },
				{ role: 'admiral', revealed: false }
			],
			isCurrentTurn: true,
			isEliminated: false,
			isYou: true
		},
		{
			id: '2',
			name: 'ShadowBlade',
			avatar: 'S',
			coins: 7,
			cards: [
				{ role: 'mercenary', revealed: true },
				{ role: 'courtesan', revealed: false }
			],
			isCurrentTurn: false,
			isEliminated: false,
			isYou: false
		},
		{
			id: '3',
			name: 'CryptoKing',
			avatar: 'C',
			coins: 2,
			cards: [
				{ role: 'consul', revealed: false },
				{ role: 'baron', revealed: false }
			],
			isCurrentTurn: false,
			isEliminated: false,
			isYou: false
		},
		{
			id: '4',
			name: 'PokerFace99',
			avatar: 'P',
			coins: 0,
			cards: [
				{ role: 'admiral', revealed: true },
				{ role: 'mercenary', revealed: true }
			],
			isCurrentTurn: false,
			isEliminated: true,
			isYou: false
		}
	]

	const maxSeats = 6

	// Position layouts per player count
	const layouts: Record<number, string[]> = {
		2: ['bottom', 'top'],
		3: ['bottom', 'top-left', 'top-right'],
		4: ['bottom', 'top', 'top-right', 'top-left'],
		5: ['bottom', 'top', 'top-left', 'top-right', 'bottom-left'],
		6: ['bottom', 'bottom-right', 'top-right', 'top-left', 'bottom-left', 'top']
	}

	const allPositions = layouts[6]
	const playerPositions = layouts[players.length] ?? layouts[6]
	const vacantPositions = allPositions.filter(
		(p) => !playerPositions.includes(p)
	)

	interface SeatData {
		player: Player | null
		position: string
	}

	const seats: SeatData[] = [
		...players.map((p, i) => ({ player: p, position: playerPositions[i] })),
		...vacantPositions.map((pos) => ({ player: null, position: pos }))
	]

	const treasury = 28
	const deckRemaining = 7

	const gameLog = [
		{
			id: 1,
			text: 'Game started — 4/6 players joined',
			type: 'system' as const
		},
		{
			id: 2,
			text: 'ShadowBlade claimed Baron — took Tax (3 coins)',
			type: 'action' as const
		},
		{ id: 3, text: 'You took Income (1 coin)', type: 'action' as const },
		{
			id: 4,
			text: 'CryptoKing claimed Admiral — stole 2 coins from PokerFace99',
			type: 'action' as const
		},
		{
			id: 5,
			text: 'PokerFace99 challenged CryptoKing — challenge failed!',
			type: 'challenge' as const
		},
		{ id: 6, text: 'PokerFace99 lost Admiral', type: 'elimination' as const },
		{
			id: 7,
			text: 'ShadowBlade claimed Mercenary — assassinated PokerFace99',
			type: 'action' as const
		},
		{
			id: 8,
			text: 'PokerFace99 lost Mercenary — eliminated!',
			type: 'elimination' as const
		},
		{ id: 9, text: 'You took Income (1 coin)', type: 'action' as const },
		{ id: 10, text: 'Your turn — choose an action', type: 'system' as const }
	]

	interface GameAction {
		id: string
		label: string
		description: string
		cost?: number
		requires?: CardRole
		needsTarget?: boolean
	}

	const actions: GameAction[] = [
		{ id: 'income', label: 'Income', description: 'Take 1 coin' },
		{ id: 'foreign-aid', label: 'Foreign Aid', description: 'Take 2 coins' },
		{
			id: 'coup',
			label: 'Coup',
			description: 'Pay 7 coins, eliminate a card',
			cost: 7,
			needsTarget: true
		},
		{ id: 'tax', label: 'Tax', description: 'Take 3 coins', requires: 'baron' },
		{
			id: 'assassinate',
			label: 'Assassinate',
			description: 'Pay 3 coins, eliminate a card',
			cost: 3,
			requires: 'mercenary',
			needsTarget: true
		},
		{
			id: 'exchange',
			label: 'Exchange',
			description: 'Swap cards with deck',
			requires: 'consul'
		},
		{
			id: 'steal',
			label: 'Steal',
			description: 'Take 2 coins from a player',
			requires: 'admiral',
			needsTarget: true
		}
	]

	const you = $derived(players.find((p) => p.isYou)!)
</script>

<div class="game">
	<div class="table-area">
		<TableCenter {deckRemaining} {treasury} />

		<div class="players-ring">
			{#each seats as seat}
				<PlayerSeat player={seat.player} position={seat.position} />
			{/each}
		</div>

		<ActionBar {actions} coins={you.coins} />
	</div>

	<GameLog entries={gameLog} />
</div>

<style lang="scss">
	.game {
		display: grid;
		grid-template-columns: 1fr 340px;
		height: calc(100vh - 60px);
		overflow: hidden;
	}

	.table-area {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background:
			radial-gradient(
				ellipse 60% 50% at 50% 50%,
				color-mix(in srgb, var(--color-primary) 12%, transparent),
				transparent 70%
			),
			radial-gradient(
				circle at 20% 80%,
				color-mix(in srgb, #0ed3cf 6%, transparent),
				transparent 40%
			),
			radial-gradient(
				circle at 80% 20%,
				color-mix(in srgb, #f43f8c 5%, transparent),
				transparent 40%
			),
			linear-gradient(
				180deg,
				color-mix(in srgb, var(--color-bg) 97%, var(--color-primary)),
				var(--color-bg)
			);
		overflow: hidden;

		/* Subtle grid pattern overlay */
		&::before {
			content: '';
			position: absolute;
			inset: 0;
			background-image: radial-gradient(
				circle,
				color-mix(in srgb, var(--color-primary) 8%, transparent) 1px,
				transparent 1px
			);
			background-size: 40px 40px;
			pointer-events: none;
			opacity: 0.4;
		}

		&::after {
			content: '';
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
			width: 70%;
			height: 60%;
			border-radius: 50%;
			border: 1px solid
				color-mix(in srgb, var(--color-primary) 15%, var(--color-border));
			background: radial-gradient(
				ellipse at center,
				color-mix(in srgb, var(--color-primary) 5%, transparent),
				transparent 70%
			);
			box-shadow:
				inset 0 0 60px color-mix(in srgb, var(--color-primary) 8%, transparent),
				0 0 40px color-mix(in srgb, var(--color-primary) 5%, transparent);
			pointer-events: none;
		}
	}

	.players-ring {
		position: absolute;
		inset: 0;
		pointer-events: none;
		z-index: 3;
	}

	@media (max-width: 900px) {
		.game {
			grid-template-columns: 1fr;
			grid-template-rows: 1fr auto;
			height: auto;
			min-height: calc(100vh - 60px);
		}

		.table-area {
			min-height: 55vh;
		}
	}
	@media (max-width: 550px) {
	}
</style>
