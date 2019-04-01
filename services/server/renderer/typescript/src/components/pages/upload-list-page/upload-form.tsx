import * as React from "react";
import "./upload-form.scss";

export interface IUploadForm {
    uploadBucketUrl: string;
    uploadParams: { [key: string]: string };
}

interface IProps {
    form: IUploadForm;
}

const INPUT_PROPS = {
    "accept": "video/*"
};

export class UploadForm extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        if (!this.props.form) {
            return null;
        }

        return <form id={"upload-form"}
                     action={this.props.form.uploadBucketUrl}
                     method={"POST"}
                     encType={"multipart/form-data"}>
            {this.renderInputParams()}

            <input type={"file"} name={"file"} required={true} {...INPUT_PROPS} />
            <button className="select">select a video</button>
            <button className="upload" type={"submit"}>upload</button>
        </form>
    }

    private renderInputParams(): JSX.Element[] {
        return Object.keys(this.props.form.uploadParams).map((key) => {
            return <input key={key} name={key} readOnly={true} value={this.props.form.uploadParams[key]} hidden={true}/>;
        });
    }

}
