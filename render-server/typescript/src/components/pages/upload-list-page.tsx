import {IPageProps, Page} from "../page";
import * as React from "react";
import {registerComponent} from "preact-rpc";

interface IUploadInfo {
    videoId: string;
    info: {
        title: string;
        description: string;
        status: string;
    };
    thumbs: string[];
    videos: Array<{
        type: string;
        url: string;
    }>
}

interface IProps extends IPageProps {
    uploads: IUploadInfo[];
}

export class UploadListPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page id={"upload-success-page"} authPayload={this.props.authPayload}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    {this.props.uploads.map((upload: IUploadInfo): JSX.Element => {
                        if (upload.info.status !== "Complete") {
                            return <div>
                                <span>Current Status: {upload.info.status}</span>
                            </div>;
                        }
                        return <div>
                            <form>
                                <input type={"text"} value={upload.info.title}/>
                                <textarea>{upload.info.description}</textarea>
                            </form>

                            <video width={500} height={500} controls={true}>
                                <source src={upload.videos.filter((video) => video.type === "mp4")[0].url} type="video/mp4"/>
                            </video>
                        </div>
                    })}
                </div>
            </div>
        </Page>;
    }

}

registerComponent('UploadListPage', UploadListPage);
