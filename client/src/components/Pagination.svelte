<script lang="ts">
	let {
		currentPage = $bindable(),
		totalPages
	}: { currentPage: number; totalPages: number } = $props()

	const pageNumbers = $derived(() => {
		const pages: (number | '...')[] = []
		if (totalPages <= 7) {
			for (let i = 1; i <= totalPages; i++) pages.push(i)
		} else {
			pages.push(1)
			if (currentPage > 3) pages.push('...')
			const start = Math.max(2, currentPage - 1)
			const end = Math.min(totalPages - 1, currentPage + 1)
			for (let i = start; i <= end; i++) pages.push(i)
			if (currentPage < totalPages - 2) pages.push('...')
			pages.push(totalPages)
		}
		return pages
	})
</script>

<div class="pagination">
	<span class="page-info">
		Page {currentPage} of {totalPages}
	</span>

	<div class="page-buttons">
		<button
			class="page-btn nav-btn"
			disabled={currentPage === 1}
			onclick={() => (currentPage = 1)}
			aria-label="First page"
		>
			⟪
		</button>
		<button
			class="page-btn nav-btn"
			disabled={currentPage === 1}
			onclick={() => currentPage--}
			aria-label="Previous page"
		>
			‹
		</button>

		{#each pageNumbers() as page}
			{#if page === '...'}
				<span class="page-ellipsis">…</span>
			{:else}
				<button
					class="page-btn"
					class:active={page === currentPage}
					onclick={() => (currentPage = page)}
				>
					{page}
				</button>
			{/if}
		{/each}

		<button
			class="page-btn nav-btn"
			disabled={currentPage === totalPages}
			onclick={() => currentPage++}
			aria-label="Next page"
		>
			›
		</button>
		<button
			class="page-btn nav-btn"
			disabled={currentPage === totalPages}
			onclick={() => (currentPage = totalPages)}
			aria-label="Last page"
		>
			⟫
		</button>
	</div>
</div>

<style lang="scss">
	.pagination {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 20px 0;
		border-top: 1px solid var(--color-border);
		margin-top: 8px;
	}

	.page-info {
		font-size: 13px;
		color: var(--color-text-secondary);
		font-weight: 400;
	}

	.page-buttons {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.page-btn {
		min-width: 36px;
		height: 36px;
		padding: 0 10px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid transparent;
		border-radius: 8px;
		color: var(--color-text-secondary);
		font-family: inherit;
		font-size: 13px;
		font-weight: 500;
		cursor: pointer;
		transition:
			background 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease,
			transform 0.1s ease;

		&:hover:not(:disabled):not(.active) {
			background: var(--color-surface);
			border-color: var(--color-border);
			color: var(--color-text);
		}

		&:active:not(:disabled) {
			transform: scale(0.93);
		}

		&.active {
			background: var(--color-primary);
			color: #fff;
			border-color: var(--color-primary);
			box-shadow: 0 2px 8px
				color-mix(in srgb, var(--color-primary) 35%, transparent);
		}

		&:disabled {
			opacity: 0.3;
			cursor: not-allowed;
		}
	}

	.nav-btn {
		font-size: 16px;
		font-weight: 400;
	}

	.page-ellipsis {
		min-width: 36px;
		height: 36px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
		font-size: 14px;
		user-select: none;
	}

	@media (max-width: 550px) {
		.pagination {
			flex-direction: column;
			gap: 12px;
			align-items: center;
		}
	}

	@media (max-width: 460px) {
		.page-btn {
			min-width: 32px;
			height: 32px;
			font-size: 12px;
		}

		.nav-btn {
			font-size: 14px;
		}
	}
</style>
