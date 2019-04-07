import * as React from "react";
import {ThumbVideo} from "../../../media/thumb-video";
import {IMediaResponse} from "../media-list-page";
import {InputGroup} from "../../../bootstrap/input-group";
import {TextareaGroup} from "../../../bootstrap/textarea-group";
import {MediaPublishForm} from "./media-publish-form";

interface IProps {
    upload: IMediaResponse;
}

export class MediaItem extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <li className={"upload-list-item row"}>
            <div className={"col-sm-4"}>
                <ThumbVideo isLink={false} showMetadata={false} media={this.props.upload}/>
            </div>
            <div className={"col-sm-8"}>
                <div className={"form-container"}>
                    <form method={"POST"} action={"/media/update"}>
                        <input type={"hidden"} name={"video_id"} value={this.props.upload.id}/>
                        <InputGroup name={"title"}
                                    type={"text"}
                                    value={this.props.upload.title}
                                    placeholder={"title (optional)"}/>
                        <TextareaGroup name={"description"}
                                       value={this.props.upload.description}
                                       placeholder={"description (required)"}/>
                        <InputGroup name={"hashtags"}
                                    type={"text"}
                                    value={this.props.upload.hashtags.join(" ")}
                                    placeholder={"hashtags (optional)"}/>
                        <button className="btn btn-primary" type="submit">update</button>
                    </form>
                    <MediaPublishForm videoId={this.props.upload.id}
                                      canPublish={this.props.upload.can_publish}
                                      isPublished={this.props.upload.is_published}/>
                </div>
            </div>
        </li>;
    }

}