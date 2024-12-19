import { ref } from 'vue';

export const isLoggedIn = ref(false);

export const updateLoginState = () => {
  isLoggedIn.value = !!localStorage.getItem('sessionToken');
};

export const handleStorageChange = (event) => {
    if (event.key === 'sessionToken') {
      updateLoginState();
    }
  };

  export const handleLoginSuccess = (token) => {
    localStorage.setItem('sessionToken', token); 
    isLoggedIn.value = true;
  };