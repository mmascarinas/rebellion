<script lang="ts">
	type GameStatus = 'waiting' | 'active' | 'full'
	type GameVisibility = 'public' | 'private'

	export interface Game {
		id: string
		name: string
		host: string
		players: number
		maxPlayers: number
		status: GameStatus
		visibility: GameVisibility
		createdAgo: string
	}

	let { game }: { game: Game } = $props()

	const statusLabels: Record<GameStatus, string> = {
		waiting: 'Waiting',
		active: 'In Progress',
		full: 'Full'
	}

	const joinable = $derived(
		game.status === 'waiting' && game.visibility === 'public'
	)
</script>

<div class="game-card" class:joinable>
	<div class="card-top">
		<div class="card-badges">
			<span class="status-badge {game.status}">
				<span class="status-dot"></span>
				{statusLabels[game.status]}
			</span>
			{#if game.visibility === 'private'}
				<span class="visibility-badge">
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						width="12"
						height="12"
					>
						<rect x="3" y="11" width="18" height="11" rx="2" />
						<path d="M7 11V7a5 5 0 0 1 10 0v4" />
					</svg>
					Private
				</span>
			{/if}
		</div>
		<span class="card-time">{game.createdAgo}</span>
	</div>

	<div class="card-body">
		<h3 class="game-name">{game.name}</h3>
		<p class="game-host">
			Hosted by <span class="host-name">{game.host}</span>
		</p>
	</div>

	<div class="card-bottom">
		<div class="player-count">
			<div class="player-dots">
				{#each Array(game.maxPlayers) as _, i}
					<span class="dot" class:filled={i < game.players}></span>
				{/each}
			</div>
			<span class="player-text">
				{game.players}/{game.maxPlayers}
			</span>
		</div>

		{#if joinable}
			<button class="join-btn">Join</button>
		{:else if game.status === 'active'}
			<button class="spectate-btn">Spectate</button>
		{:else if game.visibility === 'private' && game.status === 'waiting'}
			<button class="locked-btn" disabled aria-label="Private game">
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					width="14"
					height="14"
				>
					<rect x="3" y="11" width="18" height="11" rx="2" />
					<path d="M7 11V7a5 5 0 0 1 10 0v4" />
				</svg>
			</button>
		{:else}
			<span class="full-label">Full</span>
		{/if}
	</div>
</div>

<style lang="scss">
	.game-card {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 20px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 12px;
		transition:
			border-color 0.2s ease,
			transform 0.2s ease,
			box-shadow 0.2s ease;

		&:hover {
			transform: translateY(-2px);
			box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
		}

		&.joinable {
			border-color: color-mix(
				in srgb,
				var(--color-primary) 30%,
				var(--color-border)
			);

			&:hover {
				border-color: var(--color-primary);
				box-shadow: 0 8px 32px
					color-mix(in srgb, var(--color-primary) 15%, transparent);
			}
		}
	}

	.card-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.card-badges {
		display: flex;
		gap: 6px;
	}

	.status-badge {
		display: flex;
		align-items: center;
		gap: 5px;
		padding: 3px 10px;
		border-radius: 20px;
		font-size: 11px;
		font-weight: 600;
		letter-spacing: 0.04em;
		text-transform: uppercase;

		&.waiting {
			background: color-mix(in srgb, #22c55e 12%, transparent);
			color: #22c55e;

			.status-dot {
				background: #22c55e;
				animation: pulse 2s infinite;
			}
		}

		&.active {
			background: color-mix(in srgb, #f59e0b 12%, transparent);
			color: #f59e0b;

			.status-dot {
				background: #f59e0b;
			}
		}

		&.full {
			background: color-mix(in srgb, #ef4444 12%, transparent);
			color: #ef4444;

			.status-dot {
				background: #ef4444;
			}
		}
	}

	.status-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
	}

	.visibility-badge {
		display: flex;
		align-items: center;
		gap: 4px;
		padding: 3px 8px;
		border-radius: 20px;
		font-size: 11px;
		font-weight: 500;
		background: color-mix(
			in srgb,
			var(--color-text-secondary) 10%,
			transparent
		);
		color: var(--color-text-secondary);
	}

	.card-time {
		font-size: 11px;
		color: var(--color-text-secondary);
	}

	.card-body {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.game-name {
		font-size: 16px;
		font-weight: 600;
		color: var(--color-text);
		letter-spacing: 0.02em;
	}

	.game-host {
		font-size: 12px;
		color: var(--color-text-secondary);
	}

	.host-name {
		color: var(--color-primary);
		font-weight: 500;
	}

	.card-bottom {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: auto;
	}

	.player-count {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.player-dots {
		display: flex;
		gap: 4px;
	}

	.dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: var(--color-border);
		transition: background 0.2s ease;

		&.filled {
			background: var(--color-primary);
		}
	}

	.player-text {
		font-size: 12px;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.join-btn {
		padding: 6px 20px;
		background: var(--color-primary);
		color: #fff;
		border: none;
		border-radius: 6px;
		font-family: inherit;
		font-size: 12px;
		font-weight: 600;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		cursor: pointer;
		transition:
			transform 0.15s ease,
			box-shadow 0.15s ease;

		&:hover {
			transform: translateY(-1px);
			box-shadow: 0 4px 12px
				color-mix(in srgb, var(--color-primary) 40%, transparent);
		}
	}

	.spectate-btn {
		padding: 6px 16px;
		background: none;
		border: 1px solid var(--color-border);
		border-radius: 6px;
		color: var(--color-text-secondary);
		font-family: inherit;
		font-size: 12px;
		font-weight: 500;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		cursor: pointer;
		transition:
			border-color 0.15s ease,
			color 0.15s ease;

		&:hover {
			border-color: var(--color-primary);
			color: var(--color-primary);
		}
	}

	.locked-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		background: none;
		border: 1px solid var(--color-border);
		border-radius: 6px;
		color: var(--color-text-secondary);
		cursor: not-allowed;
		opacity: 0.5;
	}

	.full-label {
		font-size: 11px;
		font-weight: 600;
		letter-spacing: 0.08em;
		text-transform: uppercase;
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
</style>
