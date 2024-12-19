<template>
    <div class="login-page">
      <form @submit.prevent="handleLoginWithErrorHandling(emit)"
      :style="{ transform: `scale(${scale})` }">
      <h1 class="login-label">Login</h1>

        <input v-model="username" id="username" type="text" placeholder="Username" required />
        
        <input v-model="password" id="password" type="password" placeholder="Password" required />
        
        <p v-if="errorMessage" class="error-label" :class="{ visible: showError }">{{ errorMessage }}</p>

        <button type="submit">Login</button>
      </form>
    </div>
  </template>

<script setup>
import { defineEmits, onMounted } from 'vue';
import { useLogin, initializeLoginPage, scale } from '@/javascript/LoginScripts.js';
const { username, password, errorMessage, showError, handleLoginWithErrorHandling } = useLogin();

const emit = defineEmits(['login-success']);

onMounted(() => {
  initializeLoginPage();
});

</script>

<style scoped src="@/style/LoginPage.css"></style>