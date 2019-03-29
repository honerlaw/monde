import * as React from "react";
import {InputGroup} from "../../bootstrap/input-group";
import {TextareaGroup} from "../../bootstrap/textarea-group";
import {IUploadInfo} from "../upload-list-page";
import {UploadPublishForm} from "./upload-publish-form";

interface IProps {
    upload: IUploadInfo;
}

export class UploadItem extends React.Component<IProps, {}> {

    private getMp4Url(): string {
        return this.props.upload.videos.filter((video) => video.type === "mp4")[0].url;
    }

    public render(): JSX.Element {
        return <li className={"upload-list-item row"}>
            <div className={"col-sm-4"}>
                <video controls={true}>
                    <source src={this.getMp4Url()} type="video/mp4"/>
                </video>
            </div>
            <div className={"col-sm-8"}>
                <div className={"form-container"}>
                    <form method={"POST"} action={"/media/update"}>
                        <input type={"hidden"} name={"video_id"} value={this.props.upload.videoId}/>
                        <InputGroup name={"title"}
                                    type={"text"}
                                    value={this.props.upload.info.title}
                                    placeholder={"title (optional)"}/>
                        <TextareaGroup name={"description"}
                                       value={this.props.upload.info.description}
                                       placeholder={"description (required)"}/>
                        <InputGroup name={"hashtags"}
                                    type={"text"}
                                    value={this.props.upload.info.hashtags.join(" ")}
                                    placeholder={"hashtags (optional)"}/>
                        <button className="btn btn-primary" type="submit">update</button>
                    </form>
                    <UploadPublishForm videoId={this.props.upload.videoId}
                                       canPublish={this.props.upload.canPublish}
                                       isPublished={this.props.upload.info.published}/>
                </div>
            </div>
        </li>;
    }

}