import React, { Component } from "react";
import { BrowserRouter, Link, Route, Switch } from "react-router-dom";
import "./App.css";
import Movies from "./pages/Movies";
import TV from "./pages/TV";

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
            <Route path="/registration" component={TV} />
            <Route path="/login" component={TV} />
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
        <Link to="/">Moovie</Link>
        <Link to="/movies">Movies</Link>
        <Link to="/tv">TV Shows</Link>
        <div class="nav-right">
          <Link to="/login">Login</Link>
        </div>
      </div>
    );
  }
}

class Footer extends Component {
  render() {
    return (
      <footer>
        <p>&copy; 2021 INFO 441</p>
      </footer>
    );
  }
}

export default App;
