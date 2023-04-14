import './App.css';
import { useState } from 'react';

function App(props) {
    const [username, setUsername] = useState(props?.value ?? '');
    const [password, setPassword] = useState(props?.value ?? '');
    const [error, setError] = useState(props?.value ?? '');

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
        .then(response => {
            //todo fetch will automatic do the next api call when it gets 301 but not reload the page, solve
          if (response.ok) {
              window.location.reload();
          }
          setError(<p className={"error_message"}>Failed to authenticate</p>)
        })
        .catch(error => {
          setError(<p className={"error_message"}>Failed to authenticate: unknown error</p>)
          console.error(error)
        });
    }

    const handleKeyPress = (event) => {
        if(event.key === 'Enter'){
            login()
        }
    }

    return (
    <div className="App">
          <body>
              <div className={"login_prompt"}>
                  {error}
                  <div className={"prompt_box"}>
                      <div className={"prompt_username"}>
                        <p className={"prompt_username_label"}>username</p>
                        <input type={"text"} value={username} onInput={e => setUsername(e.target.value)}/>
                      </div>
                      <div className={"prompt_password"}>
                        <p className={"prompt_password_label"}>password</p>
                        <input type={"password"} value={password} onInput={e => setPassword(e.target.value)} onKeyDown={handleKeyPress}/>
                      </div>
                      <div className={"prompt_submit"}>
                        <button onClick={login}>
                            Submit
                        </button>
                      </div>
                  </div>
              </div>
          </body>
    </div>
    );
}

export default App;
