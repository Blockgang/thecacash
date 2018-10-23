import React, { Component } from 'react';
import PropTypes from 'prop-types';
import {Grid} from 'semantic-ui-react';


export default class Dashboard extends Component {
    render() {
        const {} = this.props;
        return (
            <div>
                <Grid stackable container>
                    <Grid.Row>
                        <Grid.Column computer={4} tablet={4} mobile={4}>
                        </Grid.Column>
                        <Grid.Column computer={11} tablet={11} mobile={11}>
                        </Grid.Column>
                    </Grid.Row>
                </Grid >
            </div>
        );
    }
}

Dashboard.propTypes = {
    dispatch: PropTypes.func.isRequired,
    isAuthenticated: PropTypes.bool.isRequired,
    errorMessage: PropTypes.string,
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
    instructionId: PropTypes.string,
    hintId: PropTypes.string,
    abstractID: PropTypes.string,
    solutionID: PropTypes.string,
    abstractMD: PropTypes.string,
    solutionMD: PropTypes.string,
    sectionReferences: PropTypes.array
};
