<script lang="ts">
	import ToggleSwitch from '../ToggleSwitch.svelte'
	import { toggleTheme, isDark } from '../../lib/theme.svelte'

	let mobileMenuOpen = $state(false)

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen
	}

	function closeAll() {
		mobileMenuOpen = false
	}

	function handleClickOutside(e: MouseEvent) {
		const target = e.target as HTMLElement
		if (!target.closest('.hamburger') && !target.closest('.mobile-menu')) {
			mobileMenuOpen = false
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<div class="mobile-nav">
	<button
		class="hamburger"
		class:active={mobileMenuOpen}
		aria-label="Menu"
		onclick={toggleMobileMenu}
	>
		<span class="bar"></span>
		<span class="bar"></span>
		<span class="bar"></span>
	</button>

	{#if mobileMenuOpen}
		<div class="mobile-menu">
			<button class="mobile-item" onclick={closeAll}>About Us</button>
			<button class="mobile-item" onclick={closeAll}>How to Play</button>
			<div class="mobile-divider"></div>
			<button class="mobile-item" onclick={closeAll}>
				<span class="mobile-icon">⚙</span>
				Settings
			</button>
			<button class="mobile-item" onclick={closeAll}>
				<span class="mobile-icon">★</span>
				Subscriptions
			</button>
			<div class="mobile-divider"></div>
			<div class="mobile-item theme-row">
				<span>Dark Mode</span>
				<ToggleSwitch on={isDark()} onToggle={toggleTheme} />
			</div>
		</div>
	{/if}
</div>

<style lang="scss">
	.mobile-nav {
		display: none;
		position: relative;
	}

	.hamburger {
		display: flex;
		flex-direction: column;
		justify-content: center;
		gap: 5px;
		width: 40px;
		height: 40px;
		background: none;
		border: 1px solid var(--color-border);
		border-radius: 8px;
		padding: 10px;
		cursor: pointer;
		transition:
			border-color 0.2s ease,
			background 0.2s ease;

		.bar {
			display: block;
			width: 100%;
			height: 2px;
			background: var(--color-text);
			border-radius: 1px;
			transition:
				transform 0.3s ease,
				opacity 0.3s ease;
		}

		&.active {
			border-color: var(--color-primary);
			background: color-mix(in srgb, var(--color-primary) 8%, transparent);

			.bar:nth-child(1) {
				transform: translateY(7px) rotate(45deg);
			}
			.bar:nth-child(2) {
				opacity: 0;
			}
			.bar:nth-child(3) {
				transform: translateY(-7px) rotate(-45deg);
			}
		}

		&:hover {
			border-color: var(--color-primary);
		}
	}

	.mobile-menu {
		position: absolute;
		top: calc(100% + 12px);
		right: 0;
		min-width: 220px;
		background: color-mix(in srgb, var(--color-bg) 95%, transparent);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid var(--color-border);
		border-radius: 12px;
		padding: 8px;
		box-shadow:
			0 16px 48px rgba(0, 0, 0, 0.15),
			0 0 0 1px color-mix(in srgb, var(--color-primary) 5%, transparent);
		animation: slideIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
		z-index: 200;
	}

	.mobile-item {
		display: flex;
		align-items: center;
		gap: 10px;
		width: 100%;
		padding: 10px 14px;
		background: none;
		border: none;
		border-radius: 8px;
		color: var(--color-text);
		font-family: inherit;
		font-size: 13px;
		font-weight: 400;
		cursor: pointer;
		transition:
			background 0.15s ease,
			color 0.15s ease;

		&:hover {
			background: color-mix(in srgb, var(--color-primary) 10%, transparent);
			color: var(--color-primary);
		}
	}

	.mobile-icon {
		font-size: 14px;
		width: 20px;
		text-align: center;
		color: var(--color-text-secondary);
	}

	.mobile-divider {
		height: 1px;
		background: var(--color-border);
		margin: 4px 8px;
	}

	.theme-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		cursor: default;

		&:hover {
			background: none;
			color: var(--color-text);
		}
	}

	@media (max-width: 550px) {
		.mobile-nav {
			display: block;
		}
	}

	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateY(-8px) scale(0.96);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}
</style>
