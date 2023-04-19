import './App.css';
import { useState } from 'react';

function App(props) {
    const [username, setUsername] = useState(props?.value ?? '');
    const [password, setPassword] = useState(props?.value ?? '');
    const [error, setError] = useState(props?.value ?? '');
    const search = window.location.search;
    const params = new URLSearchParams(search);
    let redirect = "/"
    if (params.has('redirect')) {
        const redirectEncoded = params.get('redirect');
        redirect = atob(redirectEncoded);
    }

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
          if (response.ok) {
              window.location.replace(redirect);
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
        </div>
    );
}

export default App;
