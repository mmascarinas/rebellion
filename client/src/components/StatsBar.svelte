<script lang="ts">
	export interface Stat {
		value: string | number
		label: string
	}

	let { stats }: { stats: Stat[] } = $props()
</script>

<div class="stats-bar">
	{#each stats as stat, i}
		{#if i > 0}
			<div class="stat-divider"></div>
		{/if}
		<div class="stat">
			<span class="stat-value">{stat.value}</span>
			<span class="stat-label">{stat.label}</span>
		</div>
	{/each}
</div>

<style lang="scss">
	.stats-bar {
		display: flex;
		align-items: center;
		gap: 24px;
		padding: 16px 24px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 12px;
	}

	.stat {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 8px;
	}

	.stat-value {
		font-size: 20px;
		font-weight: 700;
		color: var(--color-primary);
	}

	.stat-label {
		font-size: 11px;
		font-weight: 500;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-text-secondary);
	}

	.stat-divider {
		width: 1px;
		height: 32px;
		background: var(--color-border);
	}

	@media (max-width: 640px) {
		.stat {
			flex-direction: column;
			align-items: flex-start;
			gap: 2px;
		}
	}

	@media (max-width: 550px) {
		.stats-bar {
			gap: 16px;
			padding: 12px 16px;
		}
	}

	@media (max-width: 385px) {
		.stats-bar {
			gap: 10px;
			padding: 10px 12px;
			flex-wrap: wrap;
			justify-content: center;
		}

		.stat-value {
			font-size: 14px;
		}

		.stat-label {
			font-size: 9px;
			letter-spacing: 0.04em;
		}

		.stat-divider {
			height: 24px;
		}
	}
</style>
