import { ref } from 'vue';
import axios from 'axios';
import srp from 'secure-remote-password/client'

export const scale = ref(0.2);
export function useLogin() {
  const username = ref('');
  const password = ref('');
  const errorMessage = ref('');
  const showError = ref(false);

  const handleLoginWithErrorHandling = async (emit) => {
    try {
      await handleLogin(emit);
    } catch (error) {
      errorMessage.value = "Username or password is incorrect";
      showError.value = true;
  
      setTimeout(() => {
        showError.value = false;
        errorMessage.value = '';
      }, 5000);
    }
  };

  const handleLogin = async (emit) => {
    let serverResponse = await axios.post('http://localhost:8080/api/authenticate', {
      identifier: 1,
      username: username.value,
    });

    const { identifier, salt, serverPublicEphemeral } = serverResponse.data;
    if (identifier !== 2) throw new Error('Unexpected response identifier');
    
    const clientEphemeral = srp.generateEphemeral();

    const privateValue = srp.derivePrivateKey(
      salt, 
      username.value, 
      password.value
    );

    const clientSession = srp.deriveSession(
      clientEphemeral.secret,
      serverPublicEphemeral,
      salt,
      username.value,
      privateValue
    );

    serverResponse = await axios.post('http://localhost:8080/api/authenticate', {
      identifier: 3,
      username: username.value,
      clientPublicEphemeral: clientEphemeral.public,
      clientSessionProof: clientSession.proof,
      clientSessionKey: clientSession.key,
    });

    const { serverSessionProof, sessionToken } = serverResponse.data;

    srp.verifySession(
      clientEphemeral.public, 
      clientSession, 
      serverSessionProof
    )

    if (serverResponse.data.identifier !== 4) {
      throw new Error('Session verification failed');
    }else{
      scale.value = 0;
      emit('login-success', sessionToken);
    }
  };

  return {
    username,
    password,
    errorMessage,
    showError,
    handleLoginWithErrorHandling,
    handleLogin,
  };
}

function easeInOutQuart(x) {
  return x < 0.5 ? 8 * x * x * x * x : 1 - Math.pow(-2 * x + 2, 4) / 2; // https://easings.net/en
}

export function initializeLoginPage(){
  let start = null;
  const duration = 500;

  function animate(timestamp) {
    if (!start) start = timestamp;
    const elapsed = timestamp - start;
    const progress = Math.min(elapsed / duration, 1);
    scale.value = easeInOutQuart(progress) * 1;

    if (progress < 1) {
      requestAnimationFrame(animate);
    }
  }

  requestAnimationFrame(animate);
}