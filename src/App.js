import './App.css';
import { useState } from 'react';

function App(props) {
    const [username, setUsername] = useState(props?.value ?? '');
    const [password, setPassword] = useState(props?.value ?? '');

    const login = () => {
      const body = {
          username,
          password
      }
      fetch('auth', {
          method: 'POST',
          headers: {
              Accept: 'application.json',
              'Content-Type': 'application/json'
          },
          body: JSON.stringify(body),
        })
      .then(response => response.json())
      .catch(error => console.error(error));
    }

    return (
    <div className="App">
      <header className="App-header">
      </header>
      <body>
        <p>username</p>
        <input type={"text"} value={username} onInput={e => setUsername(e.target.value)}/>
        <p>password</p>
        <input type={"password"} value={password} onInput={e => setPassword(e.target.value)}/>
        <button onClick={login}>
            Submit
        </button>
      </body>
    </div>
    );
}

export default App;
