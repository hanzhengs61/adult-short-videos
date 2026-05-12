<template>
  <div ref="feedRoot" class="relative bg-black" style="height:100vh;overflow:hidden">

    <!-- 滚动容器 -->
    <div ref="container" class="w-full h-full overflow-y-scroll snap-y snap-mandatory no-scrollbar"
         style="-webkit-overflow-scrolling:touch"
         @scroll.passive="onContainerScroll">

      <div v-for="(v, idx) in videos" :key="v.video_id"
           :data-idx="idx"
           class="relative w-full snap-start snap-always shrink-0"
           :style="{ height: screenH + 'px' }">

        <!-- 广告 Slide -->
        <template v-if="v._isAd">
          <div class="w-full h-full bg-black flex items-center justify-center">
            <div class="w-full max-w-xs px-6 text-center space-y-5">
              <p class="text-white/30 text-xs tracking-widest uppercase">Advertisement</p>
              <div class="rounded-2xl border border-white/10 bg-white/5 backdrop-blur py-10 px-6">
                <p class="text-white font-bold text-2xl mb-3">广告位</p>
                <p class="text-white/40 text-base">mitun69 · 招商合作</p>
              </div>
              <p class="text-white/25 text-sm">↑ 滑动继续 ↑</p>
            </div>
          </div>
          <div class="absolute top-4 left-3 z-10 px-2 py-0.5 rounded bg-black/50 text-white/40 text-xs">广告</div>
        </template>

        <!-- 视频 Slide -->
        <template v-else>

        <!-- 视频元素：用 poster 属性显示封面，避免 video 元素黑色背景遮挡 cover -->
        <video
            :ref="el => { if (el) videoRefs[idx] = el }"
            :poster="v.cover_url || '/placeholder.svg'"
            class="absolute inset-0 w-full h-full object-contain"
            style="background:#000"
            playsinline preload="none"
            @click="togglePlay(idx)"
            @play="syncPlayState(idx, true)"
            @pause="syncPlayState(idx, false)"
            @ended="handleEnded(idx)"
            @loadedmetadata="syncVideoTime(idx)"
        />

        <!-- 渐变遮罩 -->
        <div v-if="!cleanMode" class="absolute inset-0 pointer-events-none"
             style="background:linear-gradient(to top,rgba(0,0,0,0.72) 0%,transparent 45%,transparent 70%,rgba(0,0,0,0.25) 100%)"/>

        <!-- 暂停图标 -->
        <transition name="fade-icon">
          <div v-if="pausedIdx === idx"
               class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="w-16 h-16 rounded-full bg-black/40 flex items-center justify-center">
              <svg class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
              </svg>
            </div>
          </div>
        </transition>

        <!-- 右侧操作栏：头像+关注 / 点赞 / 评论 / 分享 -->
        <div v-if="!cleanMode" class="feed-action-rail absolute right-2.5 bottom-24 z-10">

          <!-- 头像区域，底部悬挂 + 号关注按钮 -->
          <div class="feed-author-action">
            <button class="feed-avatar-btn" type="button" @click.stop="openAuthorDrawer(v)">
              <div class="feed-avatar">
                <img v-if="v.author_avatar" :src="v.author_avatar" :alt="v.author_name"
                     class="w-full h-full object-cover" @error="e => e.target.style.display='none'"/>
                <span v-if="!v.author_avatar && v.author_name">{{ v.author_name[0]?.toUpperCase() }}</span>
              </div>
            </button>
            <button v-if="v.author_name && !followedAuthors.has(v.author_name) && v.author_name !== userStore.userInfo?.username"
                    type="button"
                    @click.stop="toggleFollow(v.author_name)"
                    class="feed-follow-btn">
              <svg class="w-3.5 h-3.5 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
              </svg>
            </button>
          </div>

          <!-- 点赞 -->
          <button @click.stop="toggleLike(v)" type="button" class="feed-action-btn">
            <svg :class="['feed-action-icon feed-heart-icon transition-all', v._liked ? 'text-[#ff2b54]' : 'text-white']"
                 fill="currentColor" viewBox="0 0 48 48" aria-hidden="true">
              <path d="M24 42.4 20.9 39.6C10 29.8 3 23.5 3 15.8 3 9.5 7.9 4.8 14.1 4.8c3.5 0 6.9 1.6 9.1 4.2 2.2-2.6 5.6-4.2 9.1-4.2 6.2 0 11.1 4.7 11.1 11 0 7.7-7 14-17.9 23.8L24 42.4Z"/>
            </svg>
            <span class="feed-action-count">{{ fmtN(v.favorite_count) }}</span>
          </button>

          <!-- 评论 -->
          <button @click.stop="openComment(v)" type="button" class="feed-action-btn">
            <svg class="feed-action-icon text-white" fill="currentColor" viewBox="0 0 48 48" aria-hidden="true">
              <path d="M24 6C13.2 6 4.5 13.3 4.5 22.3c0 5.3 3 10 7.7 13 0 2.4-.8 4.7-2.2 6.6-.4.5.1 1.2.7 1 3.4-.9 6.4-2.3 8.8-4.2 1.4.3 2.9.4 4.5.4 10.8 0 19.5-7.3 19.5-16.4S34.8 6 24 6Zm-8.1 18.1a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5Zm8.1 0a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5Zm8.1 0a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5Z"/>
            </svg>
            <span class="feed-action-count">{{ fmtN(v.comment_count) }}</span>
          </button>

          <!-- 分享 -->
          <button @click.stop="openShare(v)" type="button" class="feed-action-btn">
            <svg class="feed-action-icon feed-share-icon text-white" fill="currentColor" viewBox="0 0 48 48" aria-hidden="true">
              <path d="M36.6 21.1 25.8 10.3c-1.1-1.1-3-.3-3 1.3v6.1c-9.5.8-16 6.7-18.1 16.5-.2 1.1 1.2 1.7 1.9.8 4.1-5.2 9.2-7.7 16.2-7.9v6.1c0 1.6 1.9 2.4 3 1.3l10.8-10.8c.7-.7.7-1.9 0-2.6Z"/>
            </svg>
            <span class="feed-action-count">{{ fmtN(v.share_count) }}</span>
          </button>
        </div>

        <!-- 底部信息 -->
        <div v-if="!cleanMode" class="absolute left-3 right-20 bottom-24">
          <!-- 作者名（@格式，点击打开用户抽屉） -->
          <div v-if="v.author_name" class="mb-1">
            <button @click.stop="openAuthorDrawer(v)"
                    class="text-white text-lg font-bold drop-shadow active:opacity-70 transition-opacity">
              @{{ v.author_name }}
            </button>
          </div>
          <!-- 文案（3 行省略，不展开） -->
          <p class="text-white text-sm leading-snug drop-shadow line-clamp-3 mb-2">{{ v.title }}</p>
          <!-- 标签行，点击跳转搜索 -->
          <div v-if="getTags(v).length" class="flex gap-1.5 mb-1.5 overflow-x-auto no-scrollbar">
            <button v-for="tag in getTags(v)" :key="tag"
                    @click.stop="goTag(tag)"
                    class="text-white text-[14px] font-medium active:opacity-60 transition-opacity
                             shrink-0 px-2 py-0.5 rounded-full border border-white/35 bg-white/10 backdrop-blur-sm">
              #{{ tag.replace(/^#/, '') }}
            </button>
          </div>
          <div class="flex items-center gap-1 text-white/45 text-xs">
            <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z"/>
            </svg>
            {{ fmtN(v.play_count) }} 次播放
          </div>
        </div>

        <!-- 底部控制栏 -->
        <div class="absolute bottom-0 left-0 right-0 z-10 bg-black/90 text-white" @click.stop @touchstart.stop.passive>
          <div class="relative h-2 bg-white/20">
            <div class="absolute inset-y-0 left-0 bg-white rounded-full pointer-events-none transition-none"
                 :style="{width: getProgress(idx)+'%'}"/>
            <input type="range" min="0" max="100" step="0.1"
                   :value="getProgress(idx)"
                   @input="seekTo(idx, Number($event.target.value))"
                   class="feed-progress-range absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                   style="touch-action:none"/>
          </div>
          <div class="flex min-h-12 items-center gap-3 px-3">
            <button class="control-icon-btn shrink-0" @click="togglePlay(idx)" :aria-label="isPlaying(idx) ? '暂停' : '播放'">
              <svg v-if="isPlaying(idx)" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z"/>
              </svg>
            </button>
            <div class="text-sm tabular-nums text-white/90 shrink-0">
              {{ fmtTime(videoTimes[idx]?.current) }} / {{ fmtTime(videoTimes[idx]?.duration || v.duration) }}
            </div>
            <div class="flex-1"></div>
            <div class="flex items-center gap-1">
              <button class="control-icon-btn" @click.stop="showSettings = !showSettings" aria-label="设置">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
              </button>
              <button class="control-icon-btn" @click="toggleMuted" :aria-label="muted ? '打开声音' : '静音'">
                <svg v-if="muted" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16 9l5 5m0-5l-5 5M5 9v6h4l5 4V5L9 9H5z"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 9v6h4l5 4V5L9 9H5z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" d="M17 9.5a4 4 0 010 5M19.5 7a7 7 0 010 10"/>
                </svg>
              </button>
              <button class="control-icon-btn" @click="toggleFullscreen" :aria-label="isFullscreen ? '退出全屏' : '全屏'">
                <svg v-if="isFullscreen" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v6H3M15 3v6h6M9 21v-6H3M15 21v-6h6"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8 3H3v5M16 3h5v5M8 21H3v-5M16 21h5v-5"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
        </template>
      </div>

      <!-- 底部加载指示 -->
      <div v-if="loading || hasMore" ref="sentinel"
           class="snap-start shrink-0 flex items-center justify-center"
           :style="{ height: screenH + 'px' }">
        <div v-if="loading" class="w-8 h-8 border-2 border-white/30 border-t-white rounded-full animate-spin"/>
      </div>
    </div>

    <!-- 左上角：返回 -->
    <div v-if="!cleanMode" class="absolute top-4 left-3 z-10">
      <button @click="goBack"
              class="w-9 h-9 rounded-full bg-black/40 backdrop-blur
                     flex items-center justify-center text-white active:scale-95 transition-transform">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/>
        </svg>
      </button>
    </div>

    <!-- 设置面板遮罩 -->
    <div v-if="showSettings" class="absolute inset-0 z-[29]" @click="showSettings = false"/>
    <!-- 设置面板 -->
    <transition name="fade-icon">
      <div v-if="showSettings"
           class="absolute bottom-14 right-3 z-30 bg-black/90 backdrop-blur rounded-xl shadow-xl overflow-hidden"
           style="min-width:190px"
           @click.stop>
        <!-- 连播 -->
        <div class="flex items-center justify-between px-4 py-3 border-b border-white/10">
          <span class="text-white/90 text-base">连播</span>
          <button @click="toggleContinuousPlay"
                  :class="['relative w-10 h-6 rounded-full transition-colors shrink-0', continuousPlay ? 'bg-pink-500' : 'bg-gray-600']">
            <span class="absolute left-0 top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform"
                  :style="{ transform: continuousPlay ? 'translateX(18px)' : 'translateX(2px)' }"/>
          </button>
        </div>
        <!-- 清屏 -->
        <div class="flex items-center justify-between px-4 py-3 border-b border-white/10">
          <span class="text-white/90 text-base">清屏</span>
          <button @click="cleanMode = !cleanMode"
                  :class="['relative w-10 h-6 rounded-full transition-colors shrink-0', cleanMode ? 'bg-pink-500' : 'bg-gray-600']">
            <span class="absolute left-0 top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform"
                  :style="{ transform: cleanMode ? 'translateX(18px)' : 'translateX(2px)' }"/>
          </button>
        </div>
        <!-- 倍速 -->
        <div class="px-4 py-3">
          <p class="text-white/55 text-xs font-medium tracking-wide mb-2.5">倍速</p>
          <div class="flex gap-1.5 flex-wrap">
            <button v-for="opt in speedOptions" :key="opt.value"
                    @click="setSpeed(opt.value)"
                    :class="['px-3 py-1 rounded-lg text-sm font-medium',
                             playbackRate === opt.value ? 'bg-pink-500 text-white' : 'bg-white/10 text-white/85']">
              {{ opt.value === 1 ? '正常' : opt.value + 'x' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 评论抽屉 -->
    <transition name="slide-up">
      <div v-if="commentVideo"
           class="absolute inset-0 z-20 flex flex-col justify-end bg-black/50"
           @click.self="commentVideo = null">
        <div class="bg-bg-surface rounded-t-2xl border-t border-border max-h-[70vh] flex flex-col">
          <div class="flex items-center justify-between px-4 py-3 border-b border-border shrink-0">
            <span class="text-base font-semibold text-text-primary">评论</span>
            <button @click="commentVideo = null" class="text-base font-semibold text-text-primary">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="flex-1 overflow-y-auto">
            <CommentSection
                :video-id="commentVideo.video_id"
                @comment-count-updated="updateCommentCount"
            />
          </div>
        </div>
      </div>
    </transition>

    <!-- 用户信息抽屉：点击头像弹出 -->
    <transition name="slide-up">
      <div v-if="authorDrawer"
           class="absolute inset-0 z-20 flex flex-col justify-end bg-black/50"
           @click.self="authorDrawer = null">
        <div class="bg-bg-surface rounded-t-2xl border-t border-border flex flex-col max-h-[75vh]">
          <!-- 关闭把手 -->
          <div class="flex justify-center pt-2.5 pb-1 shrink-0">
            <div class="w-9 h-1 rounded-full bg-white/20"></div>
          </div>
          <!-- 用户信息头部 -->
          <div class="flex items-center gap-3 px-4 py-3 shrink-0">
            <div class="w-14 h-14 rounded-full overflow-hidden border-2 border-white/20 bg-white/10 flex items-center justify-center text-xl font-bold text-white shrink-0">
              <img v-if="authorDrawer.avatar" :src="authorDrawer.avatar"
                   class="w-full h-full object-cover" @error="e => e.target.style.display='none'"/>
              <span v-if="!authorDrawer.avatar">{{ authorDrawer.name?.[0] }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-text-primary font-bold text-lg truncate">@{{ authorDrawer.name }}</p>
              <p class="text-text-muted text-sm mt-0.5">{{ authorDrawer.videoCount != null ? authorDrawer.videoCount + ' 个视频' : '' }}</p>
            </div>
            <!-- 非本人才显示关注按钮 -->
            <button v-if="authorDrawer.name !== userStore.userInfo?.username"
                    @click="toggleFollow(authorDrawer.name)"
                    :class="['px-5 py-2 rounded-full font-semibold text-base transition-all active:scale-95 shrink-0',
                             followedAuthors.has(authorDrawer.name)
                               ? 'bg-white/10 text-text-secondary border border-white/15'
                               : 'bg-primary text-white']">
              {{ followedAuthors.has(authorDrawer.name) ? '已关注' : '+ 关注' }}
            </button>
          </div>
          <!-- 视频网格 -->
          <div class="flex-1 overflow-y-auto px-3 pb-5">
            <div v-if="authorDrawer.videosLoading" class="grid grid-cols-3 gap-1.5 pt-1">
              <div v-for="i in 6" :key="i" class="aspect-[9/16] rounded-lg bg-white/5 animate-pulse"></div>
            </div>
            <div v-else-if="authorDrawer.videos?.length" class="grid grid-cols-3 gap-1.5 pt-1">
              <button v-for="v in authorDrawer.videos" :key="v.video_id"
                      @click="goToAuthorVideo(v)"
                      class="relative aspect-[9/16] rounded-lg overflow-hidden bg-black block">
                <img :src="v.cover_url" class="w-full h-full object-cover"/>
                <div class="absolute bottom-1 left-1 flex items-center gap-0.5 text-white text-[10px]">
                  <svg class="w-2.5 h-2.5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                  {{ fmtN(v.play_count) }}
                </div>
              </button>
            </div>
            <p v-else class="text-center py-8 text-text-muted text-sm">暂无视频</p>
          </div>
        </div>
      </div>
    </transition>

    <!-- 分享面板：点击纸飞机弹出 -->
    <transition name="slide-up">
      <div v-if="shareVideo"
           class="absolute inset-0 z-20 flex flex-col justify-end bg-black/50"
           @click.self="shareVideo = null">
        <div class="bg-bg-surface rounded-t-2xl border-t border-border">
          <div class="flex justify-center pt-2.5 pb-1">
            <div class="w-9 h-1 rounded-full bg-white/20"></div>
          </div>
          <div class="flex items-center justify-between px-4 py-2 border-b border-border/50">
            <span class="text-base font-semibold text-text-primary">分享给</span>
            <button @click="shareVideo = null">
              <svg class="w-5 h-5 text-text-secondary" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <!-- 分享渠道网格 -->
          <div class="grid grid-cols-4 gap-y-4 px-6 py-5">
            <!-- 复制链接 -->
            <button @click="doShareCopy" class="flex flex-col items-center gap-2">
              <div class="w-12 h-12 rounded-full bg-white/10 flex items-center justify-center">
                <svg class="w-6 h-6 text-text-primary" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
                </svg>
              </div>
              <span class="text-text-secondary text-[13px]">复制链接</span>
            </button>
            <!-- 系统分享（仅支持 Web Share API 的设备显示） -->
            <button v-if="canNativeShare" @click="doShareNative" class="flex flex-col items-center gap-2">
              <div class="w-12 h-12 rounded-full bg-white/10 flex items-center justify-center">
                <svg class="w-6 h-6 text-text-primary" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
                </svg>
              </div>
              <span class="text-text-secondary text-[13px]">系统分享</span>
            </button>
            <!-- 微信（触发系统分享，移动端自动列出微信） -->
            <button @click="doShareNative" class="flex flex-col items-center gap-2">
              <div class="w-12 h-12 rounded-full bg-[#07C160] flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 01.213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 00.167-.054l1.903-1.114a.864.864 0 01.717-.098 10.16 10.16 0 002.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.182a1.17 1.17 0 01-1.162 1.18A1.17 1.17 0 014.623 7.17c0-.653.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.182a1.17 1.17 0 01-1.162 1.18 1.17 1.17 0 01-1.162-1.18c0-.653.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 01.598.082l1.584.926a.272.272 0 00.14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.49.49 0 01.177-.554C22.924 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-7.062-6.122zm-3.318 2.333c.535 0 .969.44.969.982a.976.976 0 01-.969.983.976.976 0 01-.969-.983c0-.542.434-.982.969-.982zm6.748 0c.535 0 .969.44.969.982a.976.976 0 01-.969.983.976.976 0 01-.969-.983c0-.542.434-.982.969-.982z"/>
                </svg>
              </div>
              <span class="text-text-secondary text-[13px]">微信</span>
            </button>
            <!-- 微博 -->
            <button @click="doShareWeibo" class="flex flex-col items-center gap-2">
              <div class="w-12 h-12 rounded-full bg-[#E6162D] flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M10.098 20.323c-3.977.391-7.414-1.406-7.672-4.02-.259-2.609 2.759-5.047 6.74-5.441 3.979-.394 7.413 1.404 7.671 4.018.259 2.6-2.759 5.049-6.74 5.443zm5.833-11.09c-.545-.165-1.01-.276-1.326-.346.105-.35.183-.704.183-1.021 0-2.209-2.502-4.001-5.59-4.001S3.609 5.657 3.609 7.866c0 2.208 2.502 4.001 5.59 4.001.326 0 .661-.023.989-.068.081.357.259.672.535.915-1.138.044-2.221-.074-3.158-.326-1.108-.302-2.151-.867-3.007-1.666a.596.596 0 00-.849.016.616.616 0 00.016.861c1.001.944 2.209 1.597 3.479 1.952.975.266 2.012.404 3.066.404 1.048 0 2.075-.136 3.039-.395 1.271-.346 2.479-.992 3.484-1.924a.616.616 0 00.021-.867.596.596 0 00-.849-.011c-.521.49-1.106.901-1.734 1.23.235-.493.363-1.037.363-1.609 0-.146-.01-.289-.03-.431.38.088.748.195 1.045.281zm2.465-3.07c-.488-1.391-1.908-2.293-3.469-2.19l-.006-.016c-.184.015-.34.13-.415.301-.076.169-.048.365.075.507.123.143.316.203.499.155l.006.017c.957-.063 1.857.508 2.178 1.404.32.896-.024 1.881-.832 2.378l.006.016c-.163.083-.258.256-.244.44.014.184.135.342.309.403.175.061.368.012.494-.123l-.006-.016c1.212-.733 1.729-2.189 1.405-3.276zm2.364-1.036c-.978-2.788-3.82-4.598-6.945-4.388l-.01-.029c-.185.014-.342.128-.418.3-.076.171-.049.368.075.511.124.143.317.204.5.156l.01.029c2.596-.172 5.035 1.353 5.84 3.691.806 2.337-.127 4.847-2.198 6.079l.01.029c-.163.084-.258.257-.244.441.014.185.136.343.31.404.174.061.367.012.494-.123l-.009-.028c2.553-1.518 3.621-4.573 2.585-7.072z"/>
                </svg>
              </div>
              <span class="text-text-secondary text-[13px]">微博</span>
            </button>
          </div>
          <!-- 复制成功提示 -->
          <transition name="fade-icon">
            <div v-if="copyToast"
                 class="absolute bottom-32 left-1/2 -translate-x-1/2 bg-black/80 text-white text-xs px-4 py-2 rounded-full pointer-events-none whitespace-nowrap">
              链接已复制
            </div>
          </transition>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import Hls from 'hls.js'
import { videoApi, searchApi, favoriteApi, playApi, followApi } from '@/api'
import { useUserStore } from '@/stores/user'
import CommentSection from '@/components/common/CommentSection.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()


const feedRoot = ref(null)
const container = ref(null)
const videos = ref([])
const videoRefs = reactive([])   // 用 reactive 避免 ref[idx] 赋值到包装对象上
const pausedIdx = ref(-1)
const currentPlayingIdx = ref(-1)
const commentVideo = ref(null)
// 用户信息抽屉：{ name, avatar }
const authorDrawer = ref(null)
const feedBackStack = ref([])
// 分享面板当前视频
const shareVideo = ref(null)
// 复制链接成功提示
const copyToast = ref(false)
// 是否支持原生分享
const canNativeShare = typeof navigator !== 'undefined' && !!navigator.share
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
// visualViewport.height 排除了 iOS Safari 动态地址栏，避免 snap 位置错位
const screenH = ref(window.visualViewport?.height ?? window.innerHeight)
const startIdRaw = route.query.id == null ? '' : String(route.query.id)
const sentinel = ref(null)
const cleanMode = ref(false)
const continuousPlay = ref(true)
const muted = ref(false)
const isFullscreen = ref(false)
const speedOptions = [
  { value: 0.75, label: '0.75x 慢速' },
  { value: 1,    label: '1x 正常' },
  { value: 1.5,  label: '1.5x 快速' },
  { value: 2,    label: '2x 极速' },
]
const playbackRate = ref(1)
const showSettings = ref(false)
const playStates = reactive({})
const followedAuthors = reactive(new Set())
const checkedLikeIds = new Set()

// 取视频标签：优先用 v.tags 数组，否则从标题解析 #tag
function getTags(v) {
  if (Array.isArray(v.tags) && v.tags.length) return v.tags
  const matches = (v.title || '').match(/#[\w一-鿿]+/g)
  return matches ? matches.slice(0, 5) : []
}

// 点击标签跳转搜索页
function goTag(tag) {
  const q = tag.replace(/^#/, '')
  router.push({ path: '/explore', query: { q } })
}

async function checkFollowStatus(authors) {
  if (!userStore.isLoggedIn) return
  const unique = [...new Set(authors.filter(Boolean))]
  await Promise.allSettled(
    unique.map(author =>
      followApi.check(author)
        .then(res => { if (res.data?.is_followed) followedAuthors.add(author) })
        .catch(() => {})
    )
  )
}

async function toggleFollow(author) {
  if (!userStore.isLoggedIn) { userStore.openAuth('login'); return }
  try {
    if (followedAuthors.has(author)) {
      await followApi.unfollow(author)
      followedAuthors.delete(author)
    } else {
      await followApi.follow(author)
      followedAuthors.add(author)
    }
  } catch {}
}

const orderByAllowed = ['created_at', 'play_count', 'comment_count', 'favorite_count']
const orderByRaw = route.query.order_by == null ? '' : String(route.query.order_by)
const orderBy = orderByAllowed.includes(orderByRaw) ? orderByRaw : 'created_at'
// 作者过滤模式：仅加载该作者视频。改为 ref 以便抽屉切换时动态生效。
const authorFilter = ref(route.query.author ? String(route.query.author) : '')
// 搜索关键词模式：从搜索结果点进来时有值，后续加载都按此关键词过滤
const keywordFilter = ref(route.query.q ? String(route.query.q) : '')
const PAGE_SIZE = 20
const AD_EVERY = 4
const PREFETCH_SLIDES = 8
// idx -> { current: number, duration: number }
const videoTimes = reactive({})

// idx -> Hls 实例
const hlsMap = {}
// 已完成初始化的 idx
const hlsReady = new Set()
// 待播放的 idx（等 MANIFEST_PARSED 后播放）
let pendingPlay = -1
let restoringFeedFrame = false
let hydrateToken = 0

// 已移除 sentinel IntersectionObserver：初始化时 sentinel 紧贴视口
// 容易导致死循环触发 fetchMore，改用 onContainerScroll 监听滚动到底加载

function fmtN(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function stopAll() {
  videoRefs.forEach(el => {
    if (el) el.pause()
  })
}

function applyVideoSettings(el) {
  if (!el) return
  el.muted = muted.value
  el.playbackRate = playbackRate.value
}

function syncVideoTime(idx) {
  const el = videoRefs[idx]
  if (!el) return
  videoTimes[idx] = {current: el.currentTime || 0, duration: el.duration || 0}
  applyVideoSettings(el)
}

function syncPlayState(idx, playing) {
  playStates[idx] = playing
  if (playing) {
    pausedIdx.value = -1
    currentPlayingIdx.value = idx
  } else if (currentPlayingIdx.value === idx) {
    playStates[idx] = false
  }
}

function initVideo(idx) {
  if (hlsReady.has(idx) || hlsMap[idx]) return
  const v = videos.value[idx]
  const src = v?.play_url
  if (!src) return
  const el = videoRefs[idx]
  if (!el) return
  applyVideoSettings(el)

  el.addEventListener('timeupdate', () => {
    videoTimes[idx] = {current: el.currentTime, duration: el.duration || 0}
  })

  if (src.includes('.m3u8') && Hls.isSupported()) {
    // Chrome / Firefox / Edge：用 hls.js
    const hls = new Hls({enableWorker: true, maxBufferLength: 20})
    hls.loadSource(src)
    hls.attachMedia(el)
    hlsMap[idx] = hls

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      hlsReady.add(idx)
      if (pendingPlay === idx) {
        applyVideoSettings(el)
        el.play().catch(() => {
          muted.value = true
          applyVideoSettings(el)
          el.play().catch(() => {
          })
        })
      }
    })

    hls.on(Hls.Events.ERROR, (_, data) => {
      if (data.fatal) hls.destroy()
    })
    return
  }

  // Safari 原生 HLS 或 MP4
  el.src = src
  hlsReady.add(idx)
  applyVideoSettings(el)
}

function playAt(idx) {
  // 广告 slide：停止所有视频，不尝试播放任何东西
  if (videos.value[idx]?._isAd) {
    stopAll()
    pausedIdx.value = -1
    currentPlayingIdx.value = -1
    return
  }
  const el = videoRefs[idx]
  if (!el) return

  // 防止 IntersectionObserver 在列表追加/重新观察时重复触发同一条视频，
  // 导致 stopAll() 后重头播放。
  if (idx === pausedIdx.value) return
  // 已是当前播放对象：不论是否已开始播放、HLS 是否已就绪，都不再重复走 stopAll+init，
  // 否则会把刚 attachMedia 但 manifest 未解析完的视频 pause 掉，
  // 进而被浏览器 autoplay 策略 reject 而无法恢复。
  if (idx === currentPlayingIdx.value && (!el.paused || hlsMap[idx])) return

  stopAll()
  pausedIdx.value = -1
  currentPlayingIdx.value = idx
  pendingPlay = idx
  initVideo(idx)

  const v = videos.value[idx]
  checkLike(v)
  // 观看视频不应要求登录；仅记录播放历史需要登录，避免 401 导致重定向。
  if (v && userStore.isLoggedIn) playApi.record(v.video_id).catch(() => {
  })

  if (hlsReady.has(idx)) {
    applyVideoSettings(el)
    el.play().catch(() => {
      muted.value = true
      applyVideoSettings(el)
      el.play().catch(() => {
      })
    })
  }
  // 否则等 MANIFEST_PARSED 触发
}

function togglePlay(idx) {
  const el = videoRefs[idx]
  if (!el) return
  if (el.paused) {
    if (!hlsReady.has(idx)) {
      // 未初始化，走完整的 playAt 流程（含 initVideo）
      playAt(idx)
      return
    }
    applyVideoSettings(el)
    el.play().catch(() => {})
    pausedIdx.value = -1
    playStates[idx] = true
    currentPlayingIdx.value = idx
  } else {
    el.pause()
    pausedIdx.value = idx
    playStates[idx] = false
  }
}

function isPlaying(idx) {
  return !!playStates[idx]
}

async function handleEnded(idx) {
  playStates[idx] = false
  // 防御：如果 video 没有真正播放完（duration 无效、或 currentTime 远未到末尾），
  // 视为加载错误而非正常结束。这种情况自动切到下一个会陷入"疯狂下滑"循环。
  const el = videoRefs[idx]
  if (el) {
    const dur = el.duration
    if (!dur || !isFinite(dur) || dur < 1 || el.currentTime < dur - 0.5) {
      pausedIdx.value = idx
      return
    }
  }
  if (!continuousPlay.value) {
    pausedIdx.value = idx
    return
  }

  const nextIdx = idx + 1
  if (nextIdx >= videos.value.length) {
    await fetchMore()
  }
  if (nextIdx >= videos.value.length) {
    pausedIdx.value = idx
    return
  }
  pausedIdx.value = -1
  if (container.value) {
    container.value.scrollTo({top: nextIdx * screenH.value, behavior: 'smooth'})
  }
  window.setTimeout(() => playAt(nextIdx), 220)
}

function getProgress(idx) {
  const t = videoTimes[idx]
  if (!t || !t.duration) return 0
  return Math.min((t.current / t.duration) * 100, 100)
}

function seekTo(idx, pct) {
  const el = videoRefs[idx]
  if (!el || !el.duration) return
  el.currentTime = (pct / 100) * el.duration
  syncVideoTime(idx)
}

function fmtTime(s) {
  if (!s || isNaN(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = String(Math.floor(s % 60)).padStart(2, '0')
  return `${m}:${sec}`
}

let itemObservers = []
// 已注册过 observer 的元素：fetchMore 后只对新增项注册，
// 避免重置 observer 触发"初始 fire"打断当前 idx 的 HLS 加载/播放
const observedEls = new WeakSet()

function observeItems() {
  const els = container.value?.querySelectorAll('[data-idx]')
  els?.forEach(el => {
    if (observedEls.has(el)) return
    observedEls.add(el)
    const idx = Number(el.dataset.idx)
    const io = new IntersectionObserver(([entry]) => {
      if (!entry.isIntersecting || entry.intersectionRatio < 0.6) return

      // 滑动过程中可能出现相邻两条同时满足交叉阈值，
      // 用 scrollTop 推断“当前应该播放的 idx”，避免回切/重头播放。
      if (!container.value) return
      const expectedRaw = (container.value.scrollTop + screenH.value * 0.5) / screenH.value
      const expectedIdx = Math.floor(expectedRaw)
      const maxIdx = Math.max(0, videos.value.length - 1)
      const safeExpectedIdx = Math.min(Math.max(0, expectedIdx), maxIdx)
      if (idx !== safeExpectedIdx) return

      maybePrefetchFromIndex(idx)
      playAt(idx)
    }, {threshold: 0.6, root: container.value})
    io.observe(el)
    itemObservers.push(io)
  })
}

function checkLike(v) {
  if (!userStore.isLoggedIn || !v || v._isAd || !v.video_id || checkedLikeIds.has(v.video_id)) return
  checkedLikeIds.add(v.video_id)
  favoriteApi.check(v.video_id)
    .then(res => { v._liked = !!(res.data?.is_favorited ?? res.data?.is_favorite) })
    .catch(() => { checkedLikeIds.delete(v.video_id) })
}

async function fetchMore() {
  if (loading.value || !hasMore.value) return
  loading.value = true
  // 记录发出请求时的页码：若 hydrateFeedAfterTarget 在等待期间更新了列表，
  // page.value 会被改变，此时本次结果已过期，直接丢弃避免错误设置 hasMore=false
  const pageSnapshot = page.value
  try {
    const params = {page: page.value, page_size: PAGE_SIZE, order_by: orderBy}
    // 关键词模式走搜索接口，否则走列表接口
    let res
    if (keywordFilter.value) {
      res = await searchApi.videos({...params, keyword: keywordFilter.value})
    } else {
      if (authorFilter.value) params.author_name = authorFilter.value
      res = await videoApi.list(params)
    }
    const rawList = res.data?.videos || []

    // hydrateFeedAfterTarget 并发修改了 page.value，说明列表已被重建，本次结果过期
    if (page.value !== pageSnapshot) {
      loading.value = false
      return
    }

    // 过滤时跳过广告项，防止广告 ID 干扰去重逻辑
    const list = rawList
        .filter(v => !videos.value.some(e => !e._isAd && e.video_id === v.video_id))
        .map(v => ({...v, _liked: false}))

    // 每 AD_EVERY 条真实视频后插入一个广告 slide
    // 用 (globalRealIdx + 1) 使第一个广告出现在第 AD_EVERY 条视频之后（第 AD_EVERY+1 个位置）
    const realCount = videos.value.filter(v => !v._isAd).length
    const withAds = []
    list.forEach((v, i) => {
      withAds.push(v)
      const globalRealIdx = realCount + i + 1  // 1-based：第几条真实视频
      if (globalRealIdx % AD_EVERY === 0) {
        withAds.push({ _isAd: true, video_id: `_ad_${globalRealIdx}` })
      }
    })

    videos.value.push(...withAds)
    page.value++
    // 注意：这里用 rawList 判断是否到底，而不是用 filter 后长度。
    // 否则当目标视频重复被过滤时，会误判”已到底”，导致无法继续滑动。
    if (rawList.length < PAGE_SIZE) hasMore.value = false
    // 防御：如果整页全是重复（dedup 后 0 条新增），且 rawList 仍为整页大小，
    // 说明分页接口卡在原地，强制停止避免 sentinel 反复触发死循环
    if (rawList.length > 0 && list.length === 0) hasMore.value = false
    await checkFollowStatus(list.map(v => v.author_name))
    await nextTick()
    observeItems()
  } catch {
  }
  loading.value = false
}

const onResize = () => {
  screenH.value = window.visualViewport?.height ?? window.innerHeight
}

function getCurrentScrollIdx() {
  if (!container.value || !screenH.value) return 0
  return Math.max(0, Math.round(container.value.scrollTop / screenH.value))
}

function maybePrefetchFromIndex(idx = getCurrentScrollIdx()) {
  if (loading.value || !hasMore.value) return
  if (videos.value.length - idx <= PREFETCH_SLIDES) fetchMore()
}

function toggleContinuousPlay() {
  continuousPlay.value = !continuousPlay.value
}

function toggleMuted() {
  muted.value = !muted.value
  videoRefs.forEach(applyVideoSettings)
}

function setSpeed(value) {
  playbackRate.value = value
  videoRefs.forEach(applyVideoSettings)
  showSettings.value = false
}

function syncFullscreenState() {
  isFullscreen.value = !!(document.fullscreenElement || document.webkitFullscreenElement)
}

async function toggleFullscreen() {
  const target = feedRoot.value || document.documentElement
  if (!isFullscreen.value) {
    const request = target.requestFullscreen || target.webkitRequestFullscreen
    if (request) await request.call(target).catch(() => {})
  } else {
    const exit = document.exitFullscreen || document.webkitExitFullscreen
    if (exit) await exit.call(document).catch(() => {})
  }
  syncFullscreenState()
}

function goBack() {
  if (authorDrawer.value) {
    authorDrawer.value = null
    return
  }
  const frame = feedBackStack.value.pop()
  if (frame) {
    restoreFeedFrame(frame)
    return
  }
  stopAll()
  if (window.history.state?.back) {
    sessionStorage.setItem('home_feed_return_pending', '1')
    router.back()
  } else {
    router.push('/')
  }
}

// 接近尾部加载：优先按 slide 索引预取，距离底部判断作为兜底
// loading/!hasMore 时 fetchMore 自身会立即 return，无需在此判断
function onContainerScroll() {
  const el = container.value
  if (!el) return
  maybePrefetchFromIndex()
  // 距底部 5 屏时提前加载，确保数据在用户到达最后一个 snap 点前已就绪
  if (el.scrollHeight - el.scrollTop - el.clientHeight < screenH.value * 5) fetchMore()
}

async function loadInitialFeed() {
  await fetchMore()
  await nextTick()
  observeItems()
  playAt(0)
}

// 重置 FeedView 内部状态：用于"原地切换"到新起点视频（抽屉跳转、URL query 切换）
// 不修改 authorFilter，由调用方按需在调用前设置
function resetFeedState() {
  hydrateToken++  // 取消正在进行的 hydrateFeedAfterTarget
  stopAll()
  Object.values(hlsMap).forEach(h => { try { h.destroy() } catch {} })
  for (const k of Object.keys(hlsMap)) delete hlsMap[k]
  hlsReady.clear()
  itemObservers.forEach(o => o.disconnect())
  itemObservers = []
  // observedEls 是 WeakSet，DOM 销毁后旧元素自动 GC
  Object.keys(videoTimes).forEach(k => delete videoTimes[k])
  Object.keys(playStates).forEach(k => delete playStates[k])
  videoRefs.length = 0
  pausedIdx.value = -1
  currentPlayingIdx.value = -1
  pendingPlay = -1
  page.value = 1
  hasMore.value = true
  loading.value = false
  videos.value = []
}

function getCurrentVideo() {
  const idx = currentPlayingIdx.value >= 0
    ? currentPlayingIdx.value
    : Math.round((container.value?.scrollTop || 0) / screenH.value)
  return videos.value[idx] && !videos.value[idx]._isAd ? videos.value[idx] : null
}

function buildFeedFrameFromDrawer(drawer) {
  const currentVideo = getCurrentVideo()
  return {
    type: 'author-drawer',
    query: {...route.query},
    scrollTop: container.value?.scrollTop || 0,
    videoId: currentVideo?.video_id ? String(currentVideo.video_id) : '',
    drawer: {
      name: drawer.name,
      avatar: drawer.avatar,
      videos: drawer.videos || [],
      videosLoading: false,
      videoCount: drawer.videoCount,
    },
  }
}

async function restoreFeedFrame(frame) {
  stopAll()
  const idStr = frame.videoId || (frame.query?.id ? String(frame.query.id) : '')
  const canRestoreInPlace = !idStr || videos.value.some(v => !v._isAd && String(v.video_id) === idStr)

  if (frame.query && JSON.stringify(route.query) !== JSON.stringify(frame.query)) {
    authorFilter.value = frame.query.author ? String(frame.query.author) : ''
    restoringFeedFrame = true
    try {
      await router.replace({ path: '/feed', query: frame.query })
      await nextTick()
    } finally {
      restoringFeedFrame = false
    }
  }

  if (!canRestoreInPlace && idStr) {
    authorFilter.value = frame.query?.author ? String(frame.query.author) : ''
    resetFeedState()
    await loadFromVideoId(idStr)
  }

  await nextTick()
  const restoredIdx = idStr
    ? videos.value.findIndex(v => !v._isAd && String(v.video_id) === idStr)
    : -1
  const top = restoredIdx >= 0 ? restoredIdx * screenH.value : frame.scrollTop
  if (container.value) container.value.scrollTop = top || 0
  requestAnimationFrame(() => {
    const idx = restoredIdx >= 0 ? restoredIdx : Math.round((container.value?.scrollTop || 0) / screenH.value)
    playAt(idx)
  })
  if (frame.type === 'author-drawer' && frame.drawer) {
    authorDrawer.value = {...frame.drawer}
  }
}

// 每 AD_EVERY 条真实视频后插入一个广告 slide（目标视频之后的续集用）
function withAdsAfterTarget(list) {
  const result = []
  list.forEach((v, i) => {
    result.push(v)
    if ((i + 1) % AD_EVERY === 0) {
      result.push({ _isAd: true, video_id: `_ad_after_target_${i + 1}` })
    }
  })
  return result
}

// 后台搜索目标视频在列表里的真实位置，找到后把目标之后的视频按顺序拼到 Feed。
// 用 token 机制保证只有最新一次 loadFromVideoId 的 hydration 生效。
async function hydrateFeedAfterTarget(targetId, token) {
  const collected = []
  const maxSearchPages = 50
  let foundIdx = -1
  let searchPage = 1
  let lastRawLen = PAGE_SIZE

  while (token === hydrateToken && searchPage <= maxSearchPages && foundIdx < 0) {
    let rawList = []
    try {
      const params = { page: searchPage, page_size: PAGE_SIZE, order_by: orderBy }
      let res
      if (keywordFilter.value) {
        res = await searchApi.videos({ ...params, keyword: keywordFilter.value })
      } else {
        if (authorFilter.value) params.author_name = authorFilter.value
        res = await videoApi.list(params)
      }
      rawList = res.data?.videos || []
    } catch {
      return
    }

    lastRawLen = rawList.length
    for (const item of rawList) {
      if (!collected.some(v => v.video_id === item.video_id)) {
        collected.push({ ...item, _liked: false })
      }
    }
    foundIdx = collected.findIndex(v => String(v.video_id) === String(targetId))
    if (rawList.length < PAGE_SIZE) break
    searchPage++
  }

  // token 过期（用户已切视频）或找不到目标：回退到普通加载
  if (token !== hydrateToken || foundIdx < 0) {
    if (token === hydrateToken) fetchMore()
    return
  }

  const target = videos.value.find(v => !v._isAd && String(v.video_id) === String(targetId)) || collected[foundIdx]
  const afterTarget = collected
    .slice(foundIdx + 1)
    .filter(v => String(v.video_id) !== String(targetId))
    .map(v => ({ ...v, _liked: v._liked ?? false }))

  // 记住用户当前所在视频，替换列表后尽量还原位置，避免被拉回顶部
  const curIdx = container.value ? Math.round(container.value.scrollTop / screenH.value) : 0
  const curVideoId = (videos.value[curIdx] && !videos.value[curIdx]._isAd)
    ? String(videos.value[curIdx].video_id)
    : null

  videos.value = [target, ...withAdsAfterTarget(afterTarget)]
  page.value = searchPage + 1
  hasMore.value = lastRawLen >= PAGE_SIZE
  await checkFollowStatus(afterTarget.map(v => v.author_name))
  await nextTick()

  if (curIdx === 0) {
    // 用户还在目标视频，保持顶部
    if (container.value) container.value.scrollTop = 0
  } else if (curVideoId) {
    // 用户已下滑：找到当前视频在新列表里的位置
    const newIdx = videos.value.findIndex(v => !v._isAd && String(v.video_id) === curVideoId)
    if (container.value) {
      // 找到则跳到新位置；找不到（该视频在 target 之前，不在新列表里）则保持像素位置
      container.value.scrollTop = (newIdx >= 0 ? newIdx : curIdx) * screenH.value
    }
  }
  observeItems()
}

// 以单个视频为起点初始化 Feed：
//   1) 一次 detail API 拿到目标视频，立刻作为 position 0 显示并播放
//   2) 后台搜索目标在列表里的真实位置，拼接目标之后的视频（保证顺序正确）
//   3) 拉不到目标时回退到普通流
async function loadFromVideoId(targetId) {
  const token = ++hydrateToken
  let target = null
  try {
    const res = await videoApi.detail(targetId)
    target = res.data?.video || res.data || null
  } catch {}

  if (!target || !target.video_id) {
    await loadInitialFeed()
    return
  }

  target._liked = false
  videos.value = [target]
  await checkFollowStatus([target.author_name])
  await nextTick()
  if (container.value) container.value.scrollTop = 0
  observeItems()
  playAt(0)
  // 后台找到目标的真实排序位置，把目标之后的视频按顺序续上
  hydrateFeedAfterTarget(targetId, token)
}

onMounted(async () => {
  if (startIdRaw) {
    await loadFromVideoId(startIdRaw)
  } else {
    await loadInitialFeed()
  }

  window.addEventListener('resize', onResize)
  window.visualViewport?.addEventListener('resize', onResize)
  document.addEventListener('fullscreenchange', syncFullscreenState)
  document.addEventListener('webkitfullscreenchange', syncFullscreenState)
})

// 监听 URL query 变化（抽屉切视频会 router.replace，浏览器前进后退也会触发）
// 同实例原地重置内容，避免重新挂载组件 + 不增加历史栈
watch(() => route.query.id, async (newId) => {
  if (restoringFeedFrame) return
  if (!newId) return
  const idStr = String(newId)
  // 已在当前列表里：原地滚动定位（避免无谓 reset）
  const existingIdx = videos.value.findIndex(v => !v._isAd && String(v.video_id) === idStr)
  if (existingIdx >= 0) {
    if (container.value) container.value.scrollTop = existingIdx * screenH.value
    requestAnimationFrame(() => playAt(existingIdx))
    return
  }
  // 不在列表 → 切作者过滤 + 重置 + 以新视频为起点重新加载
  authorFilter.value = route.query.author ? String(route.query.author) : ''
  resetFeedState()
  await loadFromVideoId(idStr)
})

onUnmounted(() => {
  stopAll()
  Object.values(hlsMap).forEach(h => h.destroy())
  itemObservers.forEach(o => o.disconnect())
  window.removeEventListener('resize', onResize)
  window.visualViewport?.removeEventListener('resize', onResize)
  document.removeEventListener('fullscreenchange', syncFullscreenState)
  document.removeEventListener('webkitfullscreenchange', syncFullscreenState)
})

// 点赞（乐观更新：先翻转状态，API 失败时回滚）
async function toggleLike(v) {
  if (!userStore.isLoggedIn) {
    userStore.openAuth('login')
    return
  }
  const wasLiked = v._liked
  const prevCount = v.favorite_count || 0
  v._liked = !wasLiked
  v.favorite_count = wasLiked ? Math.max(0, prevCount - 1) : prevCount + 1
  try {
    if (wasLiked) {
      await favoriteApi.remove(v.video_id)
    } else {
      await favoriteApi.add(v.video_id)
    }
  } catch {
    v._liked = wasLiked
    v.favorite_count = prevCount
  }
}

function openComment(v) {
  const idx = videos.value.findIndex(item => item.video_id === v.video_id)
  if (idx >= 0) {
    initVideo(idx)
  }
  commentVideo.value = v
}

function updateCommentCount(count) {
  if (!commentVideo.value) return
  commentVideo.value.comment_count = count
  const item = videos.value.find(v => v.video_id === commentVideo.value.video_id)
  if (item) item.comment_count = count
}

async function openAuthorDrawer(v) {
  if (!v.author_name) return
  authorDrawer.value = { name: v.author_name, avatar: v.author_avatar, videos: [], videosLoading: true, videoCount: null }
  try {
    const res = await videoApi.list({ page: 1, page_size: 30, author_name: v.author_name })
    const list = res.data?.videos || []
    authorDrawer.value.videos = list
    authorDrawer.value.videoCount = res.data?.total ?? list.length
  } catch {}
  authorDrawer.value.videosLoading = false
}

function goToAuthorVideo(v) {
  const idx = videos.value.findIndex(item => !item._isAd && item.video_id === v.video_id)
  const drawer = authorDrawer.value
  if (drawer) {
    feedBackStack.value.push(buildFeedFrameFromDrawer(drawer))
  }
  authorDrawer.value = null
  if (idx >= 0) {
    // 已在列表：瞬间定位，避免看到从中间视频一路滑过去。
    if (container.value) container.value.scrollTop = idx * screenH.value
    requestAnimationFrame(() => playAt(idx))
    return
  }
  // 不在列表：用 router.replace 替换当前 URL（不开浏览器历史条目），watch 会触发原地重置 + 拉新视频。
  // Feed 内的返回关系由 feedBackStack 管理，左上角先回作者抽屉/原视频，再回上一级页面。
  router.replace({ path: '/feed', query: { author: v.author_name, id: v.video_id, order_by: orderBy } })
}

function openShare(v) {
  shareVideo.value = v
}

async function doShareCopy() {
  if (!shareVideo.value) return
  const v = shareVideo.value
  const url = `${location.origin}/feed?id=${v.video_id}&order_by=${orderBy}`
  await navigator.clipboard.writeText(url).catch(() => {})
  copyToast.value = true
  setTimeout(() => { copyToast.value = false }, 2000)
}

function doShareNative() {
  if (!shareVideo.value) return
  const v = shareVideo.value
  const url = `${location.origin}/feed?id=${v.video_id}&order_by=${orderBy}`
  if (navigator.share) {
    navigator.share({ title: v.title, url }).catch(() => {})
  } else {
    doShareCopy()
  }
  shareVideo.value = null
}

function doShareWeibo() {
  if (!shareVideo.value) return
  const v = shareVideo.value
  const url = encodeURIComponent(`${location.origin}/feed?id=${v.video_id}&order_by=${orderBy}`)
  const title = encodeURIComponent(v.title || '')
  window.open(`https://service.weibo.com/share/share.php?url=${url}&title=${title}`, '_blank')
  shareVideo.value = null
}
</script>

<style scoped>
.fade-icon-enter-active {
  transition: opacity 0.2s;
}

.fade-icon-leave-active {
  transition: opacity 0.4s;
}

.fade-icon-enter-from, .fade-icon-leave-to {
  opacity: 0;
}

.slide-up-enter-active {
  transition: transform 0.3s ease;
}

.slide-up-leave-active {
  transition: transform 0.25s ease;
}

.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(100%);
}

.feed-action-rail {
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 15px;
  width: 54px;
}

.feed-author-action {
  align-items: center;
  display: flex;
  flex-direction: column;
  margin-bottom: 2px;
  min-height: 58px;
  position: relative;
}

.feed-avatar-btn {
  border-radius: 999px;
  display: block;
  height: 46px;
  width: 46px;
}

.feed-avatar {
  align-items: center;
  background: #fff;
  border: 2px solid #fff;
  border-radius: 999px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.22);
  color: #18c9d2;
  display: flex;
  font-size: 21px;
  font-weight: 800;
  height: 46px;
  justify-content: center;
  line-height: 1;
  overflow: hidden;
  text-transform: uppercase;
  width: 46px;
}

.feed-follow-btn {
  align-items: center;
  background: #ff2b54;
  border: 2px solid #ff2b54;
  border-radius: 999px;
  bottom: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.24);
  display: flex;
  height: 22px;
  justify-content: center;
  left: 50%;
  position: absolute;
  transform: translateX(-50%);
  transition: transform 0.14s ease;
  width: 22px;
}

.feed-follow-btn:active {
  transform: translateX(-50%) scale(0.9);
}

.feed-action-btn {
  align-items: center;
  color: #fff;
  display: flex;
  flex-direction: column;
  gap: 3px;
  line-height: 1;
  min-height: 55px;
  text-align: center;
  width: 54px;
}

.feed-action-icon {
  filter: drop-shadow(0 2px 2px rgba(0, 0, 0, 0.34));
  height: 40px;
  width: 40px;
}

.feed-heart-icon {
  height: 42px;
  width: 42px;
}

.feed-share-icon {
  height: 42px;
  width: 42px;
}

.feed-action-count {
  color: #fff;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.7));
  font-size: 13px;
  font-weight: 700;
  line-height: 16px;
  min-height: 16px;
  max-width: 54px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.control-icon-btn {
  align-items: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.92);
  display: inline-flex;
  height: 32px;
  justify-content: center;
  min-width: 32px;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.control-icon-btn:active,
.control-icon-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.control-text-btn {
  align-items: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.92);
  display: inline-flex;
  font-size: 15px;
  font-weight: 600;
  height: 32px;
  justify-content: center;
  min-width: 44px;
  padding: 0 8px;
  white-space: nowrap;
}

.control-text-btn:active,
.control-text-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.control-toggle {
  align-items: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.86);
  display: inline-flex;
  font-size: 15px;
  font-weight: 600;
  gap: 5px;
  height: 32px;
  padding: 0 6px;
  white-space: nowrap;
}

.control-toggle:active,
.control-toggle:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.toggle-dot {
  background: rgba(255, 255, 255, 0.66);
  border-radius: 999px;
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.16);
  height: 8px;
  width: 8px;
}

.control-toggle.is-on .toggle-dot {
  background: #ff2b75;
  box-shadow: 0 0 0 3px rgba(255, 43, 117, 0.24);
}

.feed-progress-range {
  appearance: none;
}

.feed-progress-range::-webkit-slider-thumb {
  appearance: none;
  height: 18px;
  width: 18px;
}

.feed-progress-range::-moz-range-thumb {
  border: 0;
  height: 18px;
  width: 18px;
}

.no-scrollbar {
  scrollbar-width: none;
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}

</style>
