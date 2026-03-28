<template lang="pug">
v-app(theme="hacker")
  v-main.app-main
    ConversionPanel
</template>

<script setup lang="ts">
import ConversionPanel from './components/ConversionPanel.vue'
</script>

<style>
.app-main {
  background: transparent;
  min-height: 100vh;
}
</style>
