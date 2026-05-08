<script lang="ts" generics="T extends string">
	let {
		value = $bindable(),
		options
	}: {
		value: T
		options: { value: T; label: string }[]
	} = $props()

	let open = $state(false)
	let triggerEl: HTMLButtonElement | undefined = $state()

	let selectedLabel = $derived(
		options.find((o) => o.value === value)?.label ?? ''
	)

	function select(opt: T) {
		value = opt
		open = false
		triggerEl?.focus()
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			open = false
			triggerEl?.focus()
		}
	}

	function handleClickOutside(e: MouseEvent) {
		if (triggerEl && !triggerEl.parentElement?.contains(e.target as Node)) {
			open = false
		}
	}
</script>

<svelte:window onclick={handleClickOutside} onkeydown={handleKeydown} />

<div class="filter-select-wrap">
	<button
		class="filter-select-trigger"
		class:open
		bind:this={triggerEl}
		onclick={() => (open = !open)}
		type="button"
	>
		<span class="trigger-label">{selectedLabel}</span>
		<svg
			class="chevron"
			width="12"
			height="12"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2.5"
		>
			<path d="M6 9l6 6 6-6" />
		</svg>
	</button>

	{#if open}
		<ul class="filter-dropdown">
			{#each options as opt}
				<li>
					<button
						class="filter-option"
						class:active={opt.value === value}
						onclick={() => select(opt.value)}
						type="button"
					>
						{opt.label}
						{#if opt.value === value}
							<svg
								class="check"
								width="14"
								height="14"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2.5"
							>
								<path d="M5 12l5 5L20 7" />
							</svg>
						{/if}
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style lang="scss">
	.filter-select-wrap {
		position: relative;
		flex: 1;
	}

	.filter-select-trigger {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 8px;
		padding: 10px 12px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 8px;
		color: var(--color-text);
		font-family: inherit;
		font-size: 13px;
		font-weight: 500;
		cursor: pointer;
		transition:
			border-color 0.2s ease,
			box-shadow 0.2s ease;

		&:focus {
			outline: none;
			border-color: var(--color-primary);
			box-shadow: 0 0 0 3px
				color-mix(in srgb, var(--color-primary) 15%, transparent);
		}

		&.open {
			border-color: var(--color-primary);
		}

		.chevron {
			color: var(--color-text-secondary);
			transition: transform 0.2s ease;
		}

		&.open .chevron {
			transform: rotate(180deg);
		}
	}

	.filter-dropdown {
		position: absolute;
		top: calc(100% + 6px);
		left: 0;
		right: 0;
		z-index: 50;
		list-style: none;
		margin: 0;
		padding: 4px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 10px;
		box-shadow:
			0 4px 16px rgba(0, 0, 0, 0.1),
			0 1px 4px rgba(0, 0, 0, 0.06);
		animation: dropIn 0.15s ease;
	}

	@keyframes dropIn {
		from {
			opacity: 0;
			transform: translateY(-4px) scale(0.97);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	.filter-option {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 9px 12px;
		border: none;
		border-radius: 7px;
		background: transparent;
		color: var(--color-text);
		font-family: inherit;
		font-size: 13px;
		font-weight: 500;
		cursor: pointer;
		transition:
			background 0.12s ease,
			color 0.12s ease;

		&:hover {
			background: color-mix(in srgb, var(--color-primary) 10%, transparent);
			color: var(--color-primary);
		}

		&.active {
			background: color-mix(in srgb, var(--color-primary) 12%, transparent);
			color: var(--color-primary);
		}

		.check {
			color: var(--color-primary);
		}
	}
</style>
