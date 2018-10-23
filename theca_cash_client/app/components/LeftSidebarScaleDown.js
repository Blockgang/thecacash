import React, {Component} from 'react';
import PropTypes from 'prop-types';
import { Sidebar, Segment, Menu, Icon, Message} from 'semantic-ui-react';
import Dashboard from './Dashboard';
import ReactTooltip from 'react-tooltip';

export default class LeftSidebarScaleDown extends Component {
    state = { visible: false, showMessage: false, message: '' };

    toggleVisibility = () => this.setState({ visible: !this.state.visible });
    toggleShowMessage = () => {
        const {message} = this.props;
        if(message === '' && this.state.showMessage) {
            this.setState({showMessage: !this.state.showMessage});
        }
        this.setState({showMessage: true});
        setTimeout(()=> {
            this.setState({showMessage: false});
        }, 4000);
    };

    handleNavigateBack = (consumerURL) => {
        window.location = consumerURL;
    };

    handleAccountClick = () => {
        // window.location = SERVER_PROTOCOL + '://' + ADDRESS_KEYCLOAK + KEYLCOAK_PROXY_PATH + '/account';
    };

    logoutUser = () => {
        // const splittedProxyPath = KEYLCOAK_PROXY_PATH.split('/');
        // window.location = SERVER_PROTOCOL + '://' + ADDRESS_KEYCLOAK + KEYLCOAK_PROXY_PATH + LOGOUT_REDIRECT_URI_PATH + '=' + SERVER_PROTOCOL + '%3A%2F%2F' + ADDRESS_KEYCLOAK + '%2F' + splittedProxyPath[1] + '%2F' + splittedProxyPath[2] + '%2F' + splittedProxyPath[3] + '%2Faccount%2F';
    };

    render() {
        const { visible } = this.state;
        const { errorContext, hasError, dispatch, message, isAuthenticated} = this.props;
        if(this.state.message !== message) {
            this.toggleShowMessage();
            this.state.message = message;
        }
        return (
            <div>
                <span onClick={this.toggleVisibility} data-tip data-for={'showMenu'} className="menuHandle">☰</span>
                {this.state.showMessage && hasError &&
                <Message negative>
                    <Message.Header>{errorContext}</Message.Header>
                    <p>{message}</p>
                </Message>
                }
                {this.state.showMessage && !hasError &&
                <Message positive>
                    <Message.Header>{message}</Message.Header>
                </Message>
                }
                <ReactTooltip id={'showMenu'} place="top" type="dark" effect="float">Show Menu</ReactTooltip>
                <Sidebar.Pushable as={Segment} style={{height: '93vh'}}>
                    <Sidebar as={Menu} animation="overlay" width="thin" visible={visible} icon="labeled" vertical inverted>
                        <Menu.Item name="home" onClick={() => this.handleNavigateBack('')}>
                            <Icon name="home" />
                            Home
                        </Menu.Item>
                        <Menu.Item name="account" onClick={() => this.handleAccountClick()}>
                            <Icon name="address card" />
                            Account
                        </Menu.Item>
                        <Menu.Item name="signOut" onClick={() => this.logoutUser()}>
                            <Icon name="sign out" />
                            Sign out
                        </Menu.Item>
                    </Sidebar>
                    <Sidebar.Pusher>
                        <Segment basic>
                            <Dashboard dispatch={dispatch} isAuthenticated={isAuthenticated}/>
                        </Segment>
                    </Sidebar.Pusher>
                </Sidebar.Pushable>
            </div>
        );
    }
}

LeftSidebarScaleDown.propTypes = {
    dispatch: PropTypes.func.isRequired,
    isAuthenticated: PropTypes.bool.isRequired,
    hasError: PropTypes.bool.isRequired,
    message: PropTypes.string,
    errorContext: PropTypes.string,
    fatChallenge: PropTypes.object,
    challengeLevels: PropTypes.array,
    challengeUsages: PropTypes.array,
    uploadedFiles: PropTypes.array,
    mediaIds: PropTypes.array,
    titleImageFile: PropTypes.object,
    titleImageId: PropTypes.string,
    challengeTypes: PropTypes.array,
    challengeList: PropTypes.array,
    languageIsoCode: PropTypes.string,
    challengeId: PropTypes.string,
    challengeName: PropTypes.string,
    goldnuggetType: PropTypes.string,
    staticGoldnuggetSecret: PropTypes.string,
    challengeTitle: PropTypes.string,
    challengeType: PropTypes.string,
    challengeLevel: PropTypes.string,
    selectedChallengeUsages: PropTypes.array,
    challengeCategories: PropTypes.array,
    challengeCategory: PropTypes.string,
    challengeKeywords: PropTypes.array,
    selectedKeywords: PropTypes.array,
    isPrivate: PropTypes.string,
    actualEditorStep: PropTypes.string,
    lastEditorStep: PropTypes.string,
    selectedChallengeCategories: PropTypes.array,
    stepArray: PropTypes.array,
    sectionCount: PropTypes.number,
    editorArray: PropTypes.array,
    sectionItems: PropTypes.array,
    sectionId: PropTypes.string,
    instructionId: PropTypes.string,
    hintId: PropTypes.string,
    abstractID: PropTypes.string,
    solutionID: PropTypes.string,
    abstractMD: PropTypes.string,
    solutionMD: PropTypes.string,
    translatedChallenge: PropTypes.object,
    targetLanguage: PropTypes.string.isRequired,
    abstractTranslation: PropTypes.string.isRequired,
    solutionTranslation: PropTypes.string.isRequired,
    titleTranslation: PropTypes.string.isRequired,
    nameTranslation: PropTypes.string.isRequired,
    sectionTranslations: PropTypes.array.isRequired,
    stepTranslations: PropTypes.array.isRequired,
    sectionReferences: PropTypes.array
};
