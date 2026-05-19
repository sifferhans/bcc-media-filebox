<script setup lang="ts">
import { useAuth } from './composables/useAuth'
import LoginGate from './components/LoginGate.vue'

const { mustChoose } = useAuth()
</script>

<template>
  <LoginGate v-if="mustChoose" />
  <router-view v-else />
</template>
