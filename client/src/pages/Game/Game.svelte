<script lang="ts">
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

	interface Seat {
		player: Player | null
		index: number
	}

	function aliveCards(player: Player): number {
		return player.cards.filter((c) => !c.revealed).length
	}

	interface LogEntry {
		id: number
		text: string
		type: 'action' | 'challenge' | 'block' | 'elimination' | 'system'
	}

	// Mock game state
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

	const maxSeats = 7

	const seats: Seat[] = Array.from({ length: maxSeats }, (_, i) => ({
		index: i,
		player: players[i] ?? null
	}))

	const treasury = 28
	const deckRemaining = 7

	const gameLog: LogEntry[] = [
		{ id: 1, text: 'Game started — 4/7 players joined', type: 'system' },
		{
			id: 2,
			text: 'ShadowBlade claimed Baron — took Tax (3 coins)',
			type: 'action'
		},
		{ id: 3, text: 'You took Income (1 coin)', type: 'action' },
		{
			id: 4,
			text: 'CryptoKing claimed Admiral — stole 2 coins from PokerFace99',
			type: 'action'
		},
		{
			id: 5,
			text: 'PokerFace99 challenged CryptoKing — challenge failed!',
			type: 'challenge'
		},
		{ id: 6, text: 'PokerFace99 lost Admiral', type: 'elimination' },
		{
			id: 7,
			text: 'ShadowBlade claimed Mercenary — assassinated PokerFace99',
			type: 'action'
		},
		{
			id: 8,
			text: 'PokerFace99 lost Mercenary — eliminated!',
			type: 'elimination'
		},
		{ id: 9, text: 'You took Income (1 coin)', type: 'action' },
		{ id: 10, text: 'Your turn — choose an action', type: 'system' }
	]

	const roleLabels: Record<CardRole, string> = {
		courtesan: 'Courtesan',
		mercenary: 'Mercenary',
		baron: 'Baron',
		consul: 'Consul',
		admiral: 'Admiral'
	}

	const roleIcons: Record<CardRole, string> = {
		courtesan: '♛',
		mercenary: '⚔',
		baron: '♜',
		consul: '⚖',
		admiral: '⚓'
	}

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
	const canAfford = (action: GameAction) =>
		!action.cost || you.coins >= action.cost
</script>

<div class="game">
	<!-- Game Table -->
	<div class="table-area">
		<!-- Center: Deck & Treasury -->
		<div class="table-center">
			<div class="deck-stack">
				<div class="deck-card"></div>
				<div class="deck-card"></div>
				<div class="deck-card top"></div>
				<span class="deck-count">{deckRemaining}</span>
			</div>
			<div class="treasury">
				<span class="treasury-icon">●</span>
				<span class="treasury-amount">{treasury}</span>
			</div>
		</div>

		<!-- Players around the table -->
		<div class="players-ring">
			{#each seats as seat}
				{#if seat.player}
					{@const player = seat.player}
					<div
						class="player-seat seat-{seat.index} health-{aliveCards(
							player
						)} {player.isCurrentTurn ? 'active-turn' : ''} {player.isEliminated
							? 'eliminated'
							: ''} {player.isYou ? 'is-you' : ''}"
					>
						<div class="player-header">
							<div class="avatar" class:avatar-dead={player.isEliminated}>
								<span class="avatar-letter">{player.avatar}</span>
								{#if player.isCurrentTurn}
									<span class="avatar-ring"></span>
								{/if}
							</div>
							<div class="player-info">
								<span class="player-name">
									{player.name}
									{#if player.isCurrentTurn}
										<span class="turn-indicator"></span>
									{/if}
								</span>
								<span class="player-coins">
									<span class="coin-icon">●</span>
									{player.coins}
								</span>
							</div>
						</div>

						<div class="player-cards">
							{#each player.cards as card}
								{#if card.revealed}
									<div class="card revealed role-{card.role}">
										<span class="card-icon">{roleIcons[card.role]}</span>
										<span class="card-label">{roleLabels[card.role]}</span>
									</div>
								{:else if player.isYou}
									<div class="card own role-{card.role}">
										<span class="card-icon">{roleIcons[card.role]}</span>
										<span class="card-label">{roleLabels[card.role]}</span>
									</div>
								{:else}
									<div class="card facedown">
										<span class="card-back-pattern">⟡</span>
									</div>
								{/if}
							{/each}
						</div>

						<div class="health-pips">
							{#each player.cards as card}
								<span
									class="pip"
									class:pip-alive={!card.revealed}
									class:pip-dead={card.revealed}
								></span>
							{/each}
						</div>
					</div>
				{:else}
					<div class="player-seat seat-{seat.index} vacant">
						<div class="player-header">
							<div class="avatar avatar-vacant">
								<span class="avatar-letter">?</span>
							</div>
							<div class="player-info">
								<span class="player-name vacant-name">Waiting...</span>
							</div>
						</div>
					</div>
				{/if}
			{/each}
		</div>

		<!-- Action Bar (below your cards) -->
		<div class="action-bar">
			{#each actions as action}
				<button
					class="action-btn"
					class:requires-role={action.requires}
					class:needs-target={action.needsTarget}
					disabled={!canAfford(action)}
				>
					<span class="action-label">{action.label}</span>
					<span class="action-desc">{action.description}</span>
					<div class="action-meta">
						{#if action.cost}
							<span class="action-cost">
								<span class="coin-icon">●</span>{action.cost}
							</span>
						{/if}
						{#if action.requires}
							<span class="action-role">{roleLabels[action.requires]}</span>
						{/if}
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Sidebar: Log -->
	<div class="sidebar">
		<!-- Game Log -->
		<div class="game-log">
			<h2 class="panel-title">Game Log</h2>
			<div class="log-entries">
				{#each gameLog as entry}
					<div class="log-entry type-{entry.type}">
						<span class="log-text">{entry.text}</span>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>

<style lang="scss">
	.game {
		display: grid;
		grid-template-columns: 1fr 340px;
		height: calc(100vh - 60px);
		overflow: hidden;
	}

	/* ===== TABLE AREA ===== */
	.table-area {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: radial-gradient(
			ellipse at center,
			color-mix(in srgb, var(--color-primary) 10%, var(--color-bg)),
			var(--color-bg) 70%
		);
		overflow: hidden;
	}

	/* Center Deck + Treasury */
	.table-center {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		z-index: 1;
	}

	.deck-stack {
		position: relative;
		width: 72px;
		height: 100px;
	}

	.deck-card {
		position: absolute;
		inset: 0;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 8px;

		&:nth-child(1) {
			transform: translate(3px, 3px) rotate(2deg);
			opacity: 0.5;
		}
		&:nth-child(2) {
			transform: translate(1.5px, 1.5px) rotate(-1deg);
			opacity: 0.7;
		}
		&.top {
			transform: none;
			background: linear-gradient(
				135deg,
				var(--color-surface),
				color-mix(in srgb, var(--color-primary) 20%, var(--color-surface))
			);
			box-shadow: 0 4px 16px
				color-mix(in srgb, var(--color-primary) 15%, rgba(0, 0, 0, 0.15));
		}
	}

	.deck-count {
		position: absolute;
		bottom: -24px;
		left: 50%;
		transform: translateX(-50%);
		font-size: 11px;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.treasury {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 6px 14px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 20px;
		margin-top: 12px;

		.treasury-icon {
			color: #ffb800;
			font-size: 12px;
			text-shadow: 0 0 6px rgba(255, 184, 0, 0.5);
		}

		.treasury-amount {
			font-size: 13px;
			font-weight: 500;
			color: var(--color-text);
		}
	}

	/* Players Ring */
	.players-ring {
		position: absolute;
		inset: 0;
		pointer-events: none;
	}

	.player-seat {
		position: absolute;
		pointer-events: all;
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 12px 16px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 12px;
		border-top: 3px solid transparent;
		transition:
			border-color 0.2s ease,
			box-shadow 0.2s ease,
			opacity 0.3s ease,
			filter 0.3s ease;

		&.health-2 {
			border-top-color: #10b981;
			box-shadow: 0 -2px 8px color-mix(in srgb, #10b981 20%, transparent);
		}

		&.health-1 {
			border-top-color: #f59e0b;
			box-shadow: 0 -2px 8px color-mix(in srgb, #f59e0b 20%, transparent);
		}

		&.health-0 {
			border-top-color: #dc2626;
		}

		&.active-turn {
			box-shadow:
				0 0 24px color-mix(in srgb, var(--color-primary) 35%, transparent),
				inset 0 0 0 1px
					color-mix(in srgb, var(--color-primary) 50%, transparent);
			border-color: color-mix(
				in srgb,
				var(--color-primary) 60%,
				var(--color-border)
			);
		}

		&.eliminated {
			opacity: 0.35;
			filter: grayscale(0.85) brightness(0.8);
		}

		&.is-you {
			box-shadow: 0 0 18px
				color-mix(in srgb, var(--color-primary) 22%, transparent);
		}

		&.vacant {
			opacity: 0.4;
			border: 1px dashed var(--color-border);
			border-top: 3px dashed var(--color-border);
			background: color-mix(in srgb, var(--color-surface) 60%, transparent);
			box-shadow: none;
		}
	}

	.avatar-vacant {
		background: var(--color-border) !important;
		opacity: 0.7;
	}

	.vacant-name {
		font-style: italic;
		opacity: 0.7;
		color: var(--color-text-secondary) !important;
	}

	/* Positions around the table */
	.seat-0 {
		bottom: 120px;
		left: 50%;
		transform: translateX(-50%);
	}
	.seat-1 {
		top: 35%;
		right: 6%;
	}
	.seat-2 {
		top: 8%;
		right: 25%;
	}
	.seat-3 {
		top: 5%;
		left: 50%;
		transform: translateX(-50%);
	}
	.seat-4 {
		top: 8%;
		left: 25%;
	}
	.seat-5 {
		top: 35%;
		left: 6%;
	}
	.seat-6 {
		bottom: 30%;
		left: 5%;
	}

	/* Player Header (Avatar + Info) */
	.player-header {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.avatar {
		position: relative;
		width: 34px;
		height: 34px;
		border-radius: 50%;
		background: linear-gradient(
			135deg,
			color-mix(in srgb, var(--color-primary) 60%, var(--color-surface)),
			color-mix(in srgb, var(--color-primary) 90%, var(--color-surface))
		);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		box-shadow: 0 3px 12px
			color-mix(in srgb, var(--color-primary) 30%, transparent);

		&.avatar-dead {
			background: var(--color-border);
			box-shadow: none;
			opacity: 0.6;
		}
	}

	.avatar-letter {
		font-size: 13px;
		font-weight: 700;
		color: #fff;
		text-transform: uppercase;
	}

	.avatar-ring {
		position: absolute;
		inset: -4px;
		border-radius: 50%;
		border: 2.5px solid var(--color-primary);
		animation: ring-spin 2.5s linear infinite;
		border-top-color: transparent;
		border-right-color: transparent;
		filter: drop-shadow(0 0 3px var(--color-primary));
	}

	@keyframes ring-spin {
		to {
			transform: rotate(360deg);
		}
	}

	.player-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.player-name {
		font-size: 12px;
		font-weight: 500;
		color: var(--color-text);
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.turn-indicator {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-primary);
		animation: pulse-dot 1.5s ease infinite;
	}

	@keyframes pulse-dot {
		0%,
		100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.5;
			transform: scale(1.4);
		}
	}

	/* Health pips */
	.health-pips {
		display: flex;
		gap: 4px;
		justify-content: center;
		margin-top: 2px;
	}

	.pip {
		width: 10px;
		height: 4px;
		border-radius: 2px;
		transition:
			background 0.3s ease,
			box-shadow 0.3s ease;

		&.pip-alive {
			background: #10b981;
			box-shadow: 0 0 6px #10b981;
		}

		&.pip-dead {
			background: color-mix(in srgb, #dc2626 40%, var(--color-border));
		}
	}

	.player-coins {
		display: flex;
		align-items: center;
		gap: 4px;
		font-size: 11px;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.coin-icon {
		color: #ffb800;
		font-size: 10px;
		text-shadow: 0 0 4px rgba(255, 184, 0, 0.4);
	}

	/* Player Cards */
	.player-cards {
		display: flex;
		gap: 6px;
	}

	.card {
		width: 48px;
		height: 68px;
		border-radius: 6px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2px;
		font-size: 10px;
		font-weight: 500;
		border: 1px solid var(--color-border);
		transition: transform 0.15s ease;

		&.facedown {
			background: linear-gradient(
				135deg,
				color-mix(in srgb, var(--color-primary) 6%, var(--color-surface)),
				color-mix(in srgb, var(--color-primary) 15%, var(--color-surface))
			);
			border-color: color-mix(
				in srgb,
				var(--color-primary) 20%,
				var(--color-border)
			);
			box-shadow: 0 2px 8px
				color-mix(in srgb, var(--color-primary) 10%, transparent);
			.card-back-pattern {
				font-size: 18px;
				color: var(--color-primary);
				opacity: 0.4;
			}
		}

		&.own {
			border-color: color-mix(
				in srgb,
				var(--color-primary) 55%,
				var(--color-border)
			);
			background: color-mix(
				in srgb,
				var(--color-primary) 12%,
				var(--color-surface)
			);
			box-shadow: 0 2px 10px
				color-mix(in srgb, var(--color-primary) 18%, transparent);
			cursor: default;

			&:hover {
				transform: translateY(-2px);
				box-shadow: 0 4px 14px
					color-mix(in srgb, var(--color-primary) 25%, transparent);
			}
		}

		&.revealed {
			opacity: 0.4;
			background: var(--color-surface);
			filter: grayscale(0.5);
			border-style: dashed;
		}
	}

	.card-icon {
		font-size: 16px;
	}

	.card-label {
		font-size: 8px;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-text-secondary);
	}

	/* Role colors via card icons */
	.role-courtesan .card-icon {
		color: #f43f8c;
		filter: drop-shadow(0 0 3px rgba(244, 63, 140, 0.4));
	}
	.role-mercenary .card-icon {
		color: #f43f5e;
		filter: drop-shadow(0 0 3px rgba(244, 63, 94, 0.4));
	}
	.role-baron .card-icon {
		color: #b44dff;
		filter: drop-shadow(0 0 3px rgba(180, 77, 255, 0.4));
	}
	.role-consul .card-icon {
		color: #4d94ff;
		filter: drop-shadow(0 0 3px rgba(77, 148, 255, 0.4));
	}
	.role-admiral .card-icon {
		color: #0ed3cf;
		filter: drop-shadow(0 0 3px rgba(14, 211, 207, 0.4));
	}

	/* ===== SIDEBAR ===== */
	.sidebar {
		display: flex;
		flex-direction: column;
		border-left: 1px solid
			color-mix(in srgb, var(--color-primary) 15%, var(--color-border));
		background: linear-gradient(
			180deg,
			color-mix(in srgb, var(--color-primary) 4%, var(--color-surface)),
			var(--color-surface) 40%
		);
		overflow: hidden;
	}

	.panel-title {
		font-size: 11px;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: var(--color-primary);
		padding: 16px 16px 12px;
		text-shadow: 0 0 8px
			color-mix(in srgb, var(--color-primary) 30%, transparent);
		border-bottom: 1px solid
			color-mix(in srgb, var(--color-primary) 12%, var(--color-border));
	}

	/* Action Bar */
	.action-bar {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		padding: 14px 20px;
		background: linear-gradient(
			to top,
			color-mix(in srgb, var(--color-surface) 98%, var(--color-primary)),
			color-mix(in srgb, var(--color-surface) 90%, var(--color-primary))
		);
		backdrop-filter: blur(12px);
		border-top: 1px solid
			color-mix(in srgb, var(--color-primary) 25%, var(--color-border));
		box-shadow: 0 -4px 24px
			color-mix(in srgb, var(--color-primary) 10%, transparent);
		z-index: 10;
	}

	.action-btn {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2px;
		padding: 8px 14px 6px;
		height: 64px;
		background: linear-gradient(
			160deg,
			var(--color-bg),
			color-mix(in srgb, var(--color-primary) 6%, var(--color-bg))
		);
		border: 1px solid
			color-mix(in srgb, var(--color-primary) 20%, var(--color-border));
		border-radius: 10px;
		cursor: pointer;
		font-family: inherit;
		min-width: 0;
		box-shadow: 0 2px 6px
			color-mix(in srgb, var(--color-primary) 8%, transparent);
		transition:
			border-color 0.15s ease,
			background 0.15s ease,
			transform 0.15s ease,
			box-shadow 0.2s ease;

		&:hover:not(:disabled) {
			border-color: var(--color-primary);
			background: linear-gradient(
				160deg,
				color-mix(in srgb, var(--color-primary) 8%, var(--color-bg)),
				color-mix(in srgb, var(--color-primary) 18%, var(--color-bg))
			);
			transform: translateY(-3px);
			box-shadow:
				0 6px 20px color-mix(in srgb, var(--color-primary) 25%, transparent),
				0 0 0 1px color-mix(in srgb, var(--color-primary) 20%, transparent);
		}

		&:active:not(:disabled) {
			transform: translateY(0) scale(0.95);
			box-shadow: 0 1px 4px
				color-mix(in srgb, var(--color-primary) 15%, transparent);
		}

		&:disabled {
			opacity: 0.3;
			cursor: not-allowed;
			box-shadow: none;
			background: var(--color-surface);
			border-color: var(--color-border);
		}
	}

	.action-label {
		font-size: 11px;
		font-weight: 700;
		color: var(--color-text);
		white-space: nowrap;
		letter-spacing: 0.03em;
		text-transform: uppercase;
	}

	.action-desc {
		font-size: 9px;
		color: var(--color-text-secondary);
		white-space: nowrap;
	}

	.action-meta {
		display: flex;
		align-items: center;
		gap: 4px;
		margin-top: 2px;
	}

	.action-cost {
		display: flex;
		align-items: center;
		gap: 2px;
		font-size: 10px;
		font-weight: 700;
		color: #ffb800;
		text-shadow: 0 0 6px rgba(255, 184, 0, 0.5);
	}

	.action-role {
		font-size: 9px;
		padding: 2px 6px;
		border-radius: 4px;
		background: color-mix(in srgb, var(--color-primary) 22%, transparent);
		color: var(--color-primary);
		font-weight: 700;
		letter-spacing: 0.02em;
	}

	/* Game Log */
	.game-log {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-height: 0;
	}

	.log-entries {
		flex: 1;
		overflow-y: auto;
		padding: 8px 12px 16px;
		display: flex;
		flex-direction: column;
		gap: 5px;
	}

	.log-entry {
		padding: 8px 12px;
		border-radius: 8px;
		font-size: 12px;
		line-height: 1.4;
		color: var(--color-text-secondary);
		border-left: 3px solid transparent;
		transition:
			background 0.15s ease,
			transform 0.1s ease;

		&:hover {
			transform: translateX(2px);
		}

		&.type-action {
			border-left-color: var(--color-primary);
			background: color-mix(in srgb, var(--color-primary) 8%, transparent);
			box-shadow: inset 0 0 12px
				color-mix(in srgb, var(--color-primary) 5%, transparent);
		}
		&.type-challenge {
			border-left-color: #ffb800;
			color: #ffb800;
			background: color-mix(in srgb, #ffb800 8%, transparent);
			box-shadow: inset 0 0 12px color-mix(in srgb, #ffb800 5%, transparent);
		}
		&.type-block {
			border-left-color: #0ed3cf;
			color: #0ed3cf;
			background: color-mix(in srgb, #0ed3cf 8%, transparent);
			box-shadow: inset 0 0 12px color-mix(in srgb, #0ed3cf 5%, transparent);
		}
		&.type-elimination {
			border-left-color: #f43f5e;
			color: #f43f5e;
			background: color-mix(in srgb, #f43f5e 8%, transparent);
			box-shadow: inset 0 0 12px color-mix(in srgb, #f43f5e 5%, transparent);
		}
		&.type-system {
			font-style: italic;
			color: var(--color-primary);
			opacity: 0.8;
			text-shadow: 0 0 6px
				color-mix(in srgb, var(--color-primary) 20%, transparent);
		}
	}

	.log-text {
		display: block;
		font-weight: 500;
	}

	/* ===== RESPONSIVE ===== */
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

		.sidebar {
			border-left: none;
			border-top: 1px solid var(--color-border);
			max-height: 40vh;
		}

		.action-bar {
			gap: 4px;
			padding: 10px 12px;
			flex-wrap: wrap;
		}
	}

	@media (max-width: 550px) {
		.player-seat {
			padding: 8px 10px;
		}

		.card {
			width: 40px;
			height: 56px;
		}

		.card-icon {
			font-size: 14px;
		}

		.card-label {
			font-size: 7px;
		}

		.seat-0 {
			bottom: 3%;
		}
		.seat-1 {
			right: 3%;
		}
		.seat-2 {
			top: 5%;
		}
		.seat-3 {
			left: 3%;
		}
	}
</style>
