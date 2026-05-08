<script lang="ts">
	import { toggleTheme } from '../../lib/theme'
	import logo from '../assets/images/rebellion.png'
	import ProfileMenu from './ProfileMenu.svelte'

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

<header class="header">
	<img src={logo} alt="Rebellion" class="logo" />

	<nav class="nav-right">
		<button class="nav-link">About Us</button>

		<ProfileMenu />
	</nav>

	<div class="mobile-nav">
		<button class="hamburger" aria-label="Menu" onclick={toggleMobileMenu}>
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				width="22"
				height="22"
			>
				{#if mobileMenuOpen}
					<path d="M6 6l12 12M6 18L18 6" />
				{:else}
					<path d="M4 6h16M4 12h16M4 18h16" />
				{/if}
			</svg>
		</button>

		{#if mobileMenuOpen}
			<div class="mobile-menu">
				<button class="dropdown-item" onclick={closeAll}>About Us</button>
				<div class="dropdown-divider"></div>
				<button class="dropdown-item" onclick={closeAll}>
					<span class="dropdown-icon">⚙</span>
					Settings
				</button>
				<button class="dropdown-item" onclick={closeAll}>
					<span class="dropdown-icon">★</span>
					Subscriptions
				</button>
				<div class="dropdown-divider"></div>
				<button
					class="dropdown-item"
					onclick={() => {
						toggleTheme()
						closeAll()
					}}
				>
					<span class="dropdown-icon">◑</span>
					Toggle Theme
				</button>
			</div>
		{/if}
	</div>
</header>

<style lang="scss">
	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		border-bottom: 1px solid var(--color-border);

		.logo {
			outline: 1px solid red;
			height: clamp(90px, 10vw, 128px);
			width: auto;

			:global([data-theme='dark']) & {
				filter: invert(1) hue-rotate(180deg);
			}
		}
	}

	.nav-right {
		display: flex;
		align-items: center;
		gap: 24px;
	}

	.nav-link {
		background: none;
		border: none;
		color: var(--color-text-secondary);
		font-family: inherit;
		font-size: 13px;
		font-weight: 500;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		cursor: pointer;
		transition: color 0.15s ease;
		padding: 6px 0;

		&:hover {
			color: var(--color-text);
		}
	}

	.mobile-nav {
		display: none;
		position: relative;
	}

	.hamburger {
		background: none;
		border: 1px solid var(--color-border);
		border-radius: 6px;
		padding: 8px;
		color: var(--color-text);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: border-color 0.15s ease;

		&:hover {
			border-color: var(--color-primary);
		}
	}

	.mobile-menu {
		position: absolute;
		top: calc(100% + 8px);
		right: 0;
		min-width: 200px;
		background: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: 8px;
		padding: 6px 0;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
		animation: fadeIn 0.15s ease;
		z-index: 200;
	}

	@media (max-width: 550px) {
		.nav-right {
			display: none;
		}

		.mobile-nav {
			display: block;
		}
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
