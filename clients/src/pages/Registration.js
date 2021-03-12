import { Button, Container } from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import { makeStyles } from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import { Share } from "@material-ui/icons";
import React, { useState } from "react";
import { Link, useHistory } from "react-router-dom";

const useStyles = makeStyles((theme) => ({
  root: {
    "& .MuiTextField-root": {
      margin: theme.spacing(1),
      width: "15 ch",
    },
  },
  login: {
    paddingLeft: "20px",
    paddingRight: "20px",
    padding: "20px",
    textAlign: "center",
    borderRadius: 3,
    backgroundColor: "#E5E5E5",
    /* Center vertically and horizontally */
  },
  forms: {
    textAlign: "left",
    justifyContent: "space-evenly",
    paddingTop: 10,
    paddingBottom: 10,
  },
  birthday: {
    textAlign: "left",
    width: "8rem",
    justifyContent: "space-around",
  },
}));

export function Registration() {
  const classes = useStyles();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [username, setUsername] = useState("");
  const [team, setTeam] = useState("");
  const [menuYear, setMenuYear] = useState(null);
  const [menuMonth, setMenuMonth] = useState(null);
  const [required, setRequired] = useState(false);
  const [bio, setBio] = useState(null);
  const [hobbies, setHobbies] = useState(null);
  const [position, setPosition] = useState(null);

  var today = new Date();

  var thisMonth = today.getDate() > 4 ? today.getMonth() : today.getMonth() - 1;
  var thisYear = today.getFullYear();
  const range = (start, stop, step) =>
    Array.from(
      { length: (stop - start) / step + 1 },
      (_, i) => start + i * step
    );
  const menuYears = range(thisYear, thisYear - 120, -1);

  const monthName = (mon) => {
    return [
      "January",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December",
    ][mon - 1];
  };

  let history = useHistory(); // navigation

  return (
    <Paper>
      <Container maxWidth="sm" className={classes.login}>
        <div>
          <div>
            <Share fontSize="large" />
          </div>
          <h1>Movie Ratings</h1>
        </div>
        <div>
          <h1> Welcome! </h1>
        </div>
        <div>
          <h4> Please enter your information below to create your account</h4>
        </div>
        <Grid className={classes.forms}>
          <Grid>
            <h3> Name: </h3>
            <form className={classes.root} noValidate autoComplete="off">
              <TextField
                label="First"
                variant="outlined"
                onChange={(e) => {
                  setFirstName(e.target.value);
                }}
              />
              <TextField
                label="Last"
                variant="outlined"
                onChange={(e) => {
                  setLastName(e.target.value);
                }}
              />
            </form>
          </Grid>
          <Grid>
            <h3>Email:</h3>
            <form className={classes.root} noValidate autoComplete="off">
              <TextField
                label="Email"
                variant="outlined"
                onChange={(e) => {
                  setEmail(e.target.value);
                }}
              />
            </form>
          </Grid>
          <Grid>
            <h3>Password:</h3>
            <form className={classes.root} noValidate autoComplete="off">
              <TextField
                type="password"
                label="Password"
                variant="outlined"
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
              />
            </form>
          </Grid>
          <Grid>
            <h3>UserName:</h3>
            <form className={classes.root} noValidate autoComplete="off">
              <TextField
                label="username"
                variant="outlined"
                onChange={(e) => {
                  setUsername(e.target.value);
                }}
              />
            </form>
          </Grid>
          <Grid>
            <Button
              onClick={(e) => {
                history.push("/movies");
                fetch("https://api.jaren441.me/users/", {
                  method: "POST",
                  headers: {
                    "Content-Type": "application/json",
                  },
                  body: JSON.stringify({
                    Email: email,
                    Password: password,
                    PasswordConf: password,
                    UserName: username,
                    FirstName: firstName,
                    LastName: lastName,
                  }),
                });
              }}
            >
              Register
            </Button>
          </Grid>
          <Grid>
            <p>
              Alread have an account? Sign in <Link to={"/"}>Here</Link>
            </p>
          </Grid>
        </Grid>
      </Container>
    </Paper>
  );
}

export default Registration;
