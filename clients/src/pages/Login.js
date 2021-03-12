import { Box, Button, Container } from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import { makeStyles } from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import { Share } from "@material-ui/icons";
import { React, useState } from "react";
import { connect } from "react-redux";
import { Link, useHistory } from "react-router-dom";
import { login } from "../../../tasks";

const useStyles = makeStyles({
  root: {
    background: `#E5E5E5`,
    border: 0,
    borderRadius: 3,
    boxShadow: "0 3px 5px 2px rgba(255, 105, 135, .3)",
    color: "white",
    height: `auto`,
    width: `auto`,
    padding: "0 30px",
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
  },
});

const mapStateToProps = (state) => {
  const userId = { state };
  return {
    uid: userId,
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    login: (email, password) => dispatch(login(email, password)),
  };
};

const Login = (props) => {
  const { userId } = props.uid.state;
  const classes = useStyles();
  const [login, setLogin] = useState(null);
  const [username, setUsername] = useState(null);
  const [password, setPassword] = useState(null);
  let history = useHistory(); // navigation

  return (
    <Paper>
      <Container maxWidth="sm" className={classes.login}>
        <div>
          <div>
            <Share fontSize="large" />
          </div>
          <h1>MovieRatings</h1>
        </div>
        <div>
          <Box className={classes.forms}>
            <Grid>
              <h2> Email </h2>
              <TextField
                id="outlined-basic"
                label="Email"
                variant="outlined"
                fullWidth={true}
                onChange={(e) => {
                  setUsername(e.target.value);
                }}
              />
            </Grid>
            <Grid>
              <h2> Password </h2>
              <TextField
                type="password"
                id="outlined-basic"
                label="Password"
                variant="outlined"
                fullWidth={true}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
              />
            </Grid>
          </Box>
        </div>
        <div>
          {" "}
          <div>
            <p>
              Forgot your password? Click <span> Here</span>
            </p>
            <p>
              Don't have an account? Click{" "}
              <Link to={"/registration"}>Here</Link>
            </p>
          </div>
          <Button
            variant="contained"
            style={{ background: "F39741" }}
            onClick={(e) => {
              history.push("/");
            }}
          >
            Login
          </Button>
        </div>
      </Container>
    </Paper>
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(Login);
