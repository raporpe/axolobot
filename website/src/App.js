import logo from './logo.svg';
import './App.css';
import React, {useState} from 'react';
import {Button} from 'reactstrap';
import { applyStyles } from '@popperjs/core';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
        <ClickCounter></ClickCounter>
      </header>
    </div>
  );
}

function ClickCounter() {
  const [count, setCount] = useState(0);

  return (
    <div>
      <p>You clicked {count} times</p>
      <Button color="danger" onClick={() => setCount(count + 1)}>
        Click me
      </Button>
    </div>
  );
}

export default App;
