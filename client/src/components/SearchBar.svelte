<script lang="ts">
	let {
		value = $bindable(''),
		placeholder = 'Search...'
	}: {
		value?: string
		placeholder?: string
	} = $props()
</script>

<div class="search-wrap">
	<svg
		class="search-icon"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="2"
	>
		<circle cx="11" cy="11" r="7" />
		<path d="M21 21l-4.35-4.35" />
	</svg>
	<input class="search-input" type="text" {placeholder} bind:value />
</div>

<style lang="scss">
	.search-wrap {
		position: relative;
		flex: 1;
		min-width: 200px;
	}

	.search-icon {
		position: absolute;
		top: 50%;
		left: 14px;
		transform: translateY(-50%);
		width: 16px;
		height: 16px;
		color: var(--color-text-secondary);
		pointer-events: none;
	}

	.search-input {
		width: 100%;
		padding: 10px 14px 10px 40px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 8px;
		color: var(--color-text);
		font-family: inherit;
		font-size: 13px;
		outline: none;
		transition: border-color 0.2s ease;

		&::placeholder {
			color: var(--color-text-secondary);
		}

		&:focus {
			border-color: var(--color-primary);
		}
	}
</style>
