import {IPageProps, Page} from "../page";
import * as React from "react";
import {registerComponent} from "preact-rpc";
import {InputGroup} from "../bootstrap/input-group";

interface IProps extends IPageProps {
    uploadBucketUrl: string;
    uploadParams: { [key: string]: string };
}

const INPUT_PROPS = {
    "accept": "video/*"
};

export class UploadPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page id={"upload-page"} authPayload={this.props.authPayload} error={this.props.error}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    <form action={this.props.uploadBucketUrl} method={"POST"} encType={"multipart/form-data"}>
                        {this.renderInputParams()}
                        <InputGroup name={"file"}
                                    type={"file"}
                                    inputProps={INPUT_PROPS}
                                    placeholder={"select a video to upload"}/>
                        <button className="btn btn-primary btn-block" type="submit">upload</button>
                    </form>
                </div>
            </div>

            <script type={"text/javascript"} src={"/js/upload-form.js"} />
        </Page>;
    }

    private renderInputParams(): JSX.Element[] {
        return Object.keys(this.props.uploadParams).map((key) => {
            return <input name={key} value={this.props.uploadParams[key]} hidden={true}/>;
        });
    }

}



registerComponent('UploadPage', UploadPage);
