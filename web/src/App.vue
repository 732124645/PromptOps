<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NConfigProvider,
  NMessageProvider,
  NDialogProvider,
  darkTheme,
  zhCN,
  dateZhCN,
  enUS,
  dateEnUS,
  type GlobalThemeOverrides,
} from 'naive-ui'
import AppShell from './AppShell.vue'

const { locale } = useI18n()
const naiveLocale = computed(() => (locale.value === 'en' ? enUS : zhCN))
const naiveDateLocale = computed(() => (locale.value === 'en' ? dateEnUS : dateZhCN))

// Electric-cyan cyberpunk theme layered on top of naive-ui's dark theme.
const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#22d3ee',
    primaryColorHover: '#5fe3f5',
    primaryColorPressed: '#0e9fc0',
    primaryColorSuppl: '#22d3ee',
    infoColor: '#38bdf8',
    infoColorHover: '#5fc8fa',
    infoColorPressed: '#1f9fe0',
    infoColorSuppl: '#38bdf8',
    successColor: '#34d399',
    successColorHover: '#5fe0b3',
    warningColor: '#fbbf24',
    errorColor: '#fb7185',
    borderRadius: '8px',
    borderRadiusSmall: '6px',
    fontFamily:
      "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    fontFamilyMono: "'JetBrains Mono', 'SF Mono', Menlo, Consolas, monospace",
    bodyColor: 'transparent',
    cardColor: 'rgba(13, 20, 36, 0.6)',
    modalColor: 'rgba(13, 20, 36, 0.93)',
    popoverColor: 'rgba(15, 23, 42, 0.96)',
    tableColor: 'transparent',
    inputColor: 'rgba(7, 12, 24, 0.55)',
    borderColor: 'rgba(120, 165, 205, 0.16)',
    dividerColor: 'rgba(120, 165, 205, 0.12)',
    textColorBase: '#e6eef7',
    textColor1: '#eaf1f8',
    textColor2: '#b9c7d9',
    textColor3: '#6f8198',
    placeholderColor: '#5f7188',
    hoverColor: 'rgba(34, 211, 238, 0.09)',
  },
  Layout: {
    color: 'transparent',
    headerColor: 'transparent',
    footerColor: 'transparent',
    headerBorderColor: 'transparent',
  },
  Card: {
    color: 'rgba(13, 20, 36, 0.6)',
    borderColor: 'rgba(34, 211, 238, 0.13)',
    titleTextColor: '#eaf1f8',
  },
  Tag: {
    borderRadius: '5px',
  },
  Input: {
    borderHover: '1px solid rgba(34, 211, 238, 0.45)',
    borderFocus: '1px solid #22d3ee',
    boxShadowFocus: '0 0 0 2px rgba(34, 211, 238, 0.2)',
    caretColor: '#22d3ee',
  },
  Tabs: {
    tabTextColorActiveLine: '#22d3ee',
    barColor: '#22d3ee',
  },
}
</script>

<template>
  <n-config-provider
    :theme="darkTheme"
    :theme-overrides="themeOverrides"
    :locale="naiveLocale"
    :date-locale="naiveDateLocale"
  >
    <n-message-provider>
      <n-dialog-provider>
        <app-shell />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>
