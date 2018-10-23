import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {connect} from 'react-redux';
import Navbar from '../components/Navbar';
import LeftSidebarScaleDown from '../components/LeftSidebarScaleDown';
// import {setToken, loginUser} from '../actions/actions.js';

class App extends Component {
    render() {
        const {dispatch, isAuthenticated, errorMessage} = this.props;
        return (
            <div>
                <Navbar
                    isAuthenticated={isAuthenticated}
                    errorMessage={errorMessage}
                    dispatch={dispatch}
                />
                {isAuthenticated &&
                <LeftSidebarScaleDown dispatch={dispatch}/>
                }
            </div>
        );
    }
}

App.propTypes = {
    dispatch: PropTypes.func.isRequired,
    quote: PropTypes.string,
    isAuthenticated: PropTypes.bool.isRequired,
    errorMessage: PropTypes.string
};

function mapStateToProps(state) {
    const {auth, messageHandler} = state;
    const {isAuthenticated} = auth;
    const {hasError, message, errorContext} = messageHandler;

    return {
        isAuthenticated,
        hasError,
        errorContext,
        message
    };
}

export default connect(mapStateToProps)(App);
