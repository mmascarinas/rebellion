<script lang="ts">
  import { toggleTheme } from '../lib/theme'

  let dropdownOpen = $state(false)

  function toggleDropdown() {
    dropdownOpen = !dropdownOpen
  }

  function closeDropdown() {
    dropdownOpen = false
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as HTMLElement
    if (!target.closest('.profile')) {
      dropdownOpen = false
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

<div class="profile">
  <button class="profile-btn" aria-label="Profile menu" onclick={toggleDropdown}>
    <svg class="profile-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
      <circle cx="10.2" cy="9" r="4" />
      <path d="M4 21v-1a6 6 0 0 1 12 0v1" />
    </svg>
    <svg class="chevron" class:open={dropdownOpen} viewBox="0 0 12 8" fill="none" stroke="currentColor" stroke-width="1.5">
      <path d="M1 1.5l5 5 5-5" />
    </svg>
  </button>

  {#if dropdownOpen}
    <div class="dropdown">
      <button class="dropdown-item" onclick={closeDropdown}>
        <span class="dropdown-icon">⚙</span>
        Settings
      </button>
      <button class="dropdown-item" onclick={closeDropdown}>
        <span class="dropdown-icon">★</span>
        Subscriptions
      </button>
      <div class="dropdown-divider"></div>
      <button class="dropdown-item" onclick={() => { toggleTheme(); closeDropdown() }}>
        <span class="dropdown-icon">◑</span>
        Toggle Theme
      </button>
    </div>
  {/if}
</div>

<style lang="scss">
  .profile {
    position: relative;
  }

  .profile-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: none;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    padding: 8px 12px;
    color: var(--color-text);
    cursor: pointer;
    transition: border-color 0.15s ease;

    &:hover {
      border-color: var(--color-primary);
    }
  }

  .profile-icon {
    width: 20px;
    height: 20px;
  }

  .chevron {
    width: 10px;
    height: 10px;
    transition: transform 0.2s ease;

    &.open {
      transform: rotate(180deg);
    }
  }

  .dropdown {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    min-width: 180px;
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 6px 0;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    animation: fadeIn 0.15s ease;
    z-index: 200;
  }

  .dropdown-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 10px 16px;
    background: none;
    border: none;
    color: var(--color-text);
    font-family: inherit;
    font-size: 13px;
    cursor: pointer;
    transition: background 0.1s ease;

    &:hover {
      background: var(--color-surface);
    }
  }

  .dropdown-icon {
    font-size: 14px;
    width: 18px;
    text-align: center;
    color: var(--color-text-secondary);
  }

  .dropdown-divider {
    height: 1px;
    background: var(--color-border);
    margin: 4px 0;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(-4px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
