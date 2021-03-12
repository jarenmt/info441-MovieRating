import React, { Component } from 'react';
import './App.css';
import { Router, Switch, Route, Link, Redirect, NavLink, BrowserRouter } from 'react-router-dom';
import Movies from './pages/Movies';
import TV from './pages/TV';

class App extends Component {
  render() {
    return (
      <div>
        <BrowserRouter>
          <Nav />
          <Switch>
            <Route exact path="/" />
            <Route path="/movies" component={Movies} />
            <Route path="/tv" component={TV} />
          </Switch>
        </BrowserRouter>
        <Footer />
      </div>
    );
  }
}

class Nav extends Component {
  render() {
    return (
      <div class="nav">
        <Link to="/">Moovie</ Link>
        <Link to="/movies">Movies</ Link>
        <Link to="/tv">TV Shows</ Link>
        <div class="nav-right">
          <Link to="/login">Login</ Link>
        </div>
      </div>
    )
  }
}

class Footer extends Component {
  render() {
    return (
      <footer>
        <p>&copy;  2021 INFO 441</p>
      </footer >
    )
  }
}

export default App;
