import * as React from "react";
import {InputGroup} from "../../bootstrap/input-group";

interface IProps {
    uploadBucketUrl: string;
    uploadParams: { [key: string]: string };
}

const INPUT_PROPS = {
    "accept": "video/*"
};

export class UploadForm extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <form id={"upload-form"}
                     action={this.props.uploadBucketUrl}
                     method={"POST"}
                     encType={"multipart/form-data"}>
            {this.renderInputParams()}

            <input type={"file"} name={"file"} required={true} {...INPUT_PROPS} />
            <button className="select">select a video</button>
            <button className="upload" type={"submit"}>upload</button>
        </form>
    }

    private renderInputParams(): JSX.Element[] {
        return Object.keys(this.props.uploadParams).map((key) => {
            return <input name={key} value={this.props.uploadParams[key]} hidden={true}/>;
        });
    }

}
