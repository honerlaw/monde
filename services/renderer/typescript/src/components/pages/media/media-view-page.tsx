import * as React from "react";
import "./media-view-page.scss"
import {IGlobalProps} from "../../../global-props";
import {IMediaResponse} from "./media-list-page";
import {Video} from "../../media/video";
import {CommentList} from "./media-view-page/comment-list";
import {CommentForm} from "./media-view-page/comment-form";
import {AlertMessage} from "../../bootstrap/alert-message";

export interface ICommentResponse {
    id: string;
    user_id: string;
    media_id: string;
    parent_comment_id: string;
    comment: string;
    created_at: string;
    children: ICommentResponse[];
}

interface IProps extends IGlobalProps {
    view: IMediaResponse;
    comments: ICommentResponse[];
}

export class MediaViewPage extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        return <div id={"media-view-page"}>
            <AlertMessage type={"danger"} message={this.props.error} />
            <div className={"row"}>
                <div className={"col-sm-8 offset-sm-2"}>
                    <Video media={this.props.view}/>
                    <CommentForm mediaId={this.props.view.id}/>
                    <CommentList comments={this.props.comments}/>
                </div>
            </div>
        </div>;
    }

}
