const express = require('express');
const srp = require('secure-remote-password/server');
const app = express();
app.use(express.json());

const SubServerPort = 8082;
const SubServerAuthKey = '6IWpWBBtYzjQdCBj';

function authenticateRequest(req, res, next) {
  const authHeader = req.headers['authorization'];
  if (!authHeader || authHeader !== `Bearer ${SubServerAuthKey}`) {
    return res.status(403).send('Forbidden: Invalid API key');
  }
  next();
}

app.post('/generate-server-ephemeral', authenticateRequest, (serverRequest, subServerResponse) => {
  try {
    const { verifier, salt } = serverRequest.body;

    if (!verifier || !salt) {
      return subServerResponse.status(400).send('Verifier or Salt is missing');
    }

    const serverEphemeral = srp.generateEphemeral(verifier);

    subServerResponse.json({
      serverPublicEphemeral: serverEphemeral.public,
      serverSecretEphemeral: serverEphemeral.secret,
    });

  } catch (error) {
    subServerResponse.status(500).send('Error generating server ephemeral: ' + error.message);
  }
});

app.post('/derive-server-session', authenticateRequest, (serverRequest, subServerResponse) => {
  try {
    const { serverSecretEphemeral, clientPublicEphemeral, salt, username, verifier, clientSessionProof } = serverRequest.body;

    if (!serverSecretEphemeral || !clientPublicEphemeral || !clientSessionProof || !username || !verifier || !salt) {
      return subServerResponse.status(400).send('Missing required parameters');
    }

    const serverSession = srp.deriveSession(
      serverSecretEphemeral,
      clientPublicEphemeral,
      salt,
      username,
      verifier,
      clientSessionProof
    );

    subServerResponse.json({
      serverSessionProof: serverSession.proof,
    });
    
  } catch (error) {
    subServerResponse.status(500).send('Error deriving server session: ' + error.message);
  }
});

app.listen(SubServerPort, () => {
  console.log(`JavaScript SRP server running on port ${SubServerPort}`);
});
