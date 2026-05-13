<script lang="ts">
	type CardRole = 'courtesan' | 'mercenary' | 'baron' | 'consul' | 'admiral'

	interface GameAction {
		id: string
		label: string
		description: string
		cost?: number
		requires?: CardRole
		needsTarget?: boolean
	}

	const roleLabels: Record<CardRole, string> = {
		courtesan: 'Courtesan',
		mercenary: 'Mercenary',
		baron: 'Baron',
		consul: 'Consul',
		admiral: 'Admiral'
	}

	let { actions, coins }: { actions: GameAction[]; coins: number } = $props()

	const canAfford = (action: GameAction) => !action.cost || coins >= action.cost
</script>

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
						<span class="coin-icon">🪙</span>{action.cost}
					</span>
				{/if}
				{#if action.requires}
					<span class="action-role">{roleLabels[action.requires]}</span>
				{/if}
			</div>
		</button>
	{/each}
</div>

<style lang="scss">
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

	@media (max-width: 900px) {
		.action-bar {
			gap: 4px;
			padding: 10px 12px;
			flex-wrap: wrap;
		}
	}
</style>
