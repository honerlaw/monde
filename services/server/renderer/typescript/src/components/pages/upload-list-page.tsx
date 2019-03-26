import {IPageProps, Page} from "../page";
import * as React from "react";
import {registerComponent} from "preact-rpc";
import {InputGroup} from "../bootstrap/input-group";
import {TextareaGroup} from "../bootstrap/textarea-group";

interface IUploadInfo {
    videoId: string;
    info: {
        title: string;
        description: string;
        status: string;
        hashtags: string[];
        published: boolean;
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

/**
 * @todo
 * - no video place holder (e.g. no videos exist)
 * - display thumbnail after upload so the user can see what it is
 * - add ability to publish / unpublish a video (this will be sketch at first, since we aren't going to modify the s3 bucket itself)
 * - add a button that links to the /media/upload page
 */
export class UploadListPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page id={"upload-list-page"} authPayload={this.props.authPayload}>
            <ol className={"upload-list"}>
                {this.props.uploads.map((upload: IUploadInfo): JSX.Element => {
                    if (upload.info.status !== "Complete") {
                        return this.renderPending(upload);
                    }
                    return this.renderInfo(upload);
                })}
            </ol>
        </Page>;
    }

    private renderPending(upload: IUploadInfo): JSX.Element {
        return <li className={"upload-list-item row"}>
            <div className={"col-sm-4"}>
                <div className={"placeholder"}>
                    <span>v</span>
                </div>
            </div>
            <div className={"col-sm-8 text-center"}>
                <span className={"status"}>status: {upload.info.status}</span>
            </div>
        </li>;
    }

    private renderInfo(upload: IUploadInfo): JSX.Element {
        return <li className={"upload-list-item row"}>
            <div className={"col-sm-4"}>
                <video controls={true}>
                    <source src={this.getMp4Url(upload)} type="video/mp4"/>
                </video>
            </div>
            <div className={"col-sm-8"}>
                <div className={"form-container"}>
                    <form method={"POST"} action={"/media/update"}>
                        <input type={"hidden"} name={"video_id"} value={upload.videoId}/>
                        <InputGroup name={"title"} type={"text"} value={upload.info.title} placeholder={"title"}/>
                        <TextareaGroup name={"description"} value={upload.info.description}
                                       placeholder={"description"}/>
                        <InputGroup name={"hashtags"}
                                    type={"text"}
                                    value={upload.info.hashtags.join(" ")}
                                    placeholder={"hashtags"}/>
                        <button className="btn btn-primary" type="submit">update</button>
                    </form>
                    <form method={"POST"} action={"/media/publish"} className={"publish-form"}>
                        <input type={"hidden"} name={"video_id"} value={upload.videoId}/>
                        <button className="btn btn-primary" type={"submit"}>
                            {upload.info.published ? "unpublish" : "publish"}
                        </button>
                    </form>
                </div>
            </div>
        </li>;
    }

    private getMp4Url(upload: IUploadInfo): string {
        return upload.videos.filter((video) => video.type === "mp4")[0].url;
    }

}

registerComponent('UploadListPage', UploadListPage);
