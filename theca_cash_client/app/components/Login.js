import React, {Component} from 'react';
import {Segment, Button, Grid, Header} from 'semantic-ui-react';
import PropTypes from 'prop-types';
import ValidatedInputField from './includes/ValidatedInputField';
import {StringHelper} from '../helpers/StringHelper';
import {createPasswordHash} from '../helpers/CryptoHelper';

export default class Login extends Component {
    constructor(props) {
        super(props);
        this.state = {username: '', password: '', disabled: true};
    }

    areStateFieldsFilled = (username, password) => {
        return !StringHelper.isNullOrEmpty(username) && !StringHelper.isNullOrEmpty(password);
    };

    handleUsernameChange = (event: Event) => {
        if(event.target instanceof HTMLInputElement) {
            this.setState({username: event.target.value});
            if(this.areStateFieldsFilled(event.target.username, this.state.password)) {
                this.setState({disabled: false});
            } else {
                this.setState({disabled: true});
            }
        }

        return true;
    };

    handlePasswordChange = (event: Event) => {
        if(event.target instanceof HTMLInputElement) {
            this.setState({password: event.target.value});
            if(this.areStateFieldsFilled(this.state.username, event.target.value)) {
                this.setState({disabled: false});
            } else {
                this.setState({disabled: true});
            }
        }

        return true;
    };

    handleSubmit = (event) => {
        const {onLoginClick} = this.props;
        if(event.target instanceof HTMLButtonElement) {
            const creds = {};
            // const username = this.refs.username;
            // const password = this.refs.password;
            // const creds = {username: username.value.trim(), password: password.value.trim()};
            if(this.state !== undefined && this.state !== null) {
                creds.username = this.state.username;
                creds.password = createPasswordHash(this.state.password);
                onLoginClick(creds);
            }
        }
    };

    render() {
        const {errorMessage} = this.props;

        return (
            <div>
                <Grid container>
                    <Grid.Column>
                        <Segment raised>
                            <Header size="large">Login</Header>
                            <Segment>
                                <ValidatedInputField key={'username'}
                                    className={'simpleInputField'}
                                    size={'huge'}
                                    id={'username'}
                                    reference={'one'}
                                    value={this.state.username}
                                    focused
                                    onChange={this.handleUsernameChange}
                                    required
                                    placeholder={'Please insert username...'}
                                    name={'username'}
                                    ref={'username'}
                                    labelText={'Username'}/>
                            </Segment>
                            <Segment>
                                <ValidatedInputField key={'password'}
                                    className={'simpleInputField'}
                                    size={'huge'}
                                    id={'password'}
                                    reference={'one'}
                                    value={this.state.password}
                                    onChange={this.handlePasswordChange}
                                    required
                                    type={'password'}
                                    placeholder={'Password'}
                                    name={'password'}
                                    labelText={'Password'}/>
                            </Segment>
                            <Button primary fluid disabled={this.state.disabled} id={'loginSubmitButton'}
                                onClick={(event) => this.handleSubmit(event)}>Login</Button>
                        </Segment>
                    </Grid.Column>
                </Grid>
                {errorMessage &&
                <p style={{color: 'red'}}>{errorMessage}</p>
                }
            </div>
        );
    }
}

Login.propTypes = {
    onLoginClick: PropTypes.func.isRequired,
    errorMessage: PropTypes.string
};
