<script lang="ts" generics="T extends string">
	let {
		value = $bindable(),
		options
	}: {
		value: T
		options: { value: T; label: string }[]
	} = $props()
</script>

<div class="filter-group">
	{#each options as opt}
		<button
			class="filter-btn"
			class:active={value === opt.value}
			onclick={() => (value = opt.value)}
		>
			{opt.label}
		</button>
	{/each}
</div>

<style lang="scss">
	.filter-group {
		display: flex;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 8px;
		overflow: hidden;
	}

	.filter-btn {
		padding: 8px 14px;
		background: none;
		border: none;
		color: var(--color-text-secondary);
		font-family: inherit;
		font-size: 12px;
		font-weight: 500;
		letter-spacing: 0.04em;
		cursor: pointer;
		transition:
			background 0.15s ease,
			color 0.15s ease;

		&:hover {
			color: var(--color-text);
		}

		&.active {
			background: var(--color-primary);
			color: #fff;
		}
	}
</style>
