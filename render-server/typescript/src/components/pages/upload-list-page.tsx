import {IPageProps, Page} from "../page";
import * as React from "react";
import {registerComponent} from "preact-rpc";

interface IProps extends IPageProps {

}

export class UploadListPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page id={"upload-success-page"} authPayload={this.props.authPayload}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    successfully uploaded the file! or successfully did something...
                </div>
            </div>
        </Page>;
    }

}

registerComponent('UploadListPage', UploadListPage);
