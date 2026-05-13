<script lang="ts">
	import RoleCard from './RoleCard.svelte'

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

	interface Props {
		player: Player | null
		position: string
	}

	let { player, position }: Props = $props()

	function aliveCards(p: Player): number {
		return p.cards.filter((c) => !c.revealed).length
	}
</script>

{#if player}
	<div
		class="player-seat pos-{position} health-{aliveCards(
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
					<span class="coin-icon">🪙</span>
					{player.coins}
				</span>
			</div>
		</div>

		<div class="player-cards">
			{#each player.cards as card}
				{#if card.revealed}
					<RoleCard role={card.role} revealed />
				{:else if player.isYou}
					<RoleCard role={card.role} isOwn />
				{:else}
					<RoleCard role={card.role} facedown />
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
	<div class="player-seat pos-{position} vacant">
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

<style lang="scss">
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

	/* Named positions around the table */
	.pos-bottom {
		bottom: 120px;
		left: 50%;
		transform: translateX(-50%);
	}
	.pos-top {
		top: 10%;
		left: 50%;
		transform: translateX(-50%);
	}
	.pos-top-left {
		top: 20%;
		left: 12%;
	}
	.pos-top-right {
		top: 20%;
		right: 12%;
	}
	.pos-bottom-left {
		bottom: 35%;
		left: 6%;
	}
	.pos-bottom-right {
		bottom: 35%;
		right: 6%;
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

		&.avatar-vacant {
			background: var(--color-border);
			opacity: 0.7;
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

	.vacant-name {
		font-style: italic;
		opacity: 0.7;
		color: var(--color-text-secondary) !important;
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

	.player-coins {
		display: flex;
		align-items: center;
		gap: 4px;
		font-size: 11px;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.coin-icon {
		font-size: 11px;
		line-height: 1;
	}

	/* Player Cards */
	.player-cards {
		display: flex;
		gap: 6px;
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

	@media (max-width: 550px) {
		.player-seat {
			padding: 8px 10px;
		}

		.pos-bottom {
			bottom: 3%;
		}
		.pos-bottom-right {
			right: 3%;
		}
		.pos-top-right {
			top: 5%;
			right: 5%;
		}
		.pos-top-left {
			top: 5%;
			left: 5%;
		}
		.pos-bottom-left {
			left: 3%;
		}
	}
</style>
