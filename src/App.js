import './App.css';

function App() {
  const login = () => {
      fetch('login')
          .then(response => response.json())
          .catch(error => console.error(error));
  }

  return (
    <div className="App">
      <header className="App-header">
      </header>
      <body>
        <p>username</p>
        <input type={"text"}/>
        <p>password</p>
        <input type={"password"}/>
        <button onClick={login}>
            Submit
        </button>
      </body>
    </div>
  );
}

export default App;
