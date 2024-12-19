/*global google*/
<template>
  <div id="app">
    <LoginComponent v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
    <MapComponent v-else />
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import LoginComponent from './components/LoginComponent.vue';
import MapComponent from './components/MapComponent.vue';
import {handleStorageChange, isLoggedIn, handleLoginSuccess, updateLoginState} from './javascript/login-utilities.js';

onMounted(() => {
  updateLoginState()

  window.addEventListener('storage', handleStorageChange);
});

</script>
