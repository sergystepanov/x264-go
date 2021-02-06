// Package x264c implements cgo bindings for [x264](https://www.videolan.org/developers/x264.html) library.
package external

/*
#include "stdint.h"
#include "x264.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

/****************************************************************************
 * Encoder parameters
 ****************************************************************************/
/* CPU flags */

const (
	/* x86 */
	X264CpuMmx    uint = 1 << 0
	X264CpuMmx2   uint = 1 << 1 /* MMX2 aka MMXEXT aka ISSE */
	X264CpuMmxext      = X264CpuMmx2
	X264CpuSse    uint = 1 << 2
	X264CpuSse2   uint = 1 << 3
	X264CpuLzcnt  uint = 1 << 4
	X264CpuSse3   uint = 1 << 5
	X264CpuSsse3  uint = 1 << 6
	X264CpuSse4   uint = 1 << 7  /* SSE4.1 */
	X264CpuSse42  uint = 1 << 8  /* SSE4.2 */
	X264CpuAvx    uint = 1 << 9  /* Requires OS support even if YMM registers aren't used */
	X264CpuXop    uint = 1 << 10 /* AMD XOP */
	X264CpuFma4   uint = 1 << 11 /* AMD FMA4 */
	X264CpuFma3   uint = 1 << 12
	X264CpuBmi1   uint = 1 << 13
	X264CpuBmi2   uint = 1 << 14
	X264CpuAvx2   uint = 1 << 15
	X264CpuAvx512 uint = 1 << 16 /* AVX-512 {F, CD, BW, DQ, VL}, requires OS support */
	/* x86 modifiers */
	X264CpuCacheline32 uint = 1 << 17 /* avoid memory loads that span the border between two cachelines */
	X264CpuCacheline64 uint = 1 << 18 /* 32/64 is the size of a cacheline in bytes */
	X264CpuSse2IsSlow  uint = 1 << 19 /* avoid most SSE2 functions on Athlon64 */
	X264CpuSse2IsFast  uint = 1 << 20 /* a few functions are only faster on Core2 and Phenom */
	X264CpuSlowShuffle uint = 1 << 21 /* The Conroe has a slow shuffle unit (relative to overall SSE performance) */
	X264CpuStackMod4   uint = 1 << 22 /* if stack is only mod4 and not mod16 */
	X264CpuSlowAtom    uint = 1 << 23 /* The Atom is terrible: slow SSE unaligned loads, slow
	 * SIMD multiplies, slow SIMD variable shifts, slow pshufb,
	 * cacheline split penalties -- gather everything here that
	 * isn't shared by other CPUs to avoid making half a dozen
	 * new SLOW flags. */
	X264CpuSlowPshufb  uint = 1 << 24 /* such as on the Intel Atom */
	X264CpuSlowPalignr uint = 1 << 25 /* such as on the AMD Bobcat */

	/* PowerPC */
	X264CpuAltivec uint = 0x0000001

	/* ARM and AArch64 */
	X264CpuArmv6       uint = 0x0000001
	X264CpuNeon        uint = 0x0000002 /* ARM NEON */
	X264CpuFastNeonMrc uint = 0x0000004 /* Transfer from NEON to ARM register is fast (Cortex-A9) */
	X264CpuArmv8       uint = 0x0000008

	/* MIPS */
	X264CpuMsa uint = 0x0000001 /* MIPS MSA */

	/* Analyse flags */
	X264AnalyseI4x4      uint = 0x0001 /* Analyse i4x4 */
	X264AnalyseI8x8      uint = 0x0002 /* Analyse i8x8 (requires 8x8 transform) */
	X264AnalysePsub16x16 uint = 0x0010 /* Analyse p16x8, p8x16 and p8x8 */
	X264AnalysePsub8x8   uint = 0x0020 /* Analyse p8x4, p4x8, p4x4 */
	X264AnalyseBsub16x16 uint = 0x0100 /* Analyse b16x8, b8x16 and b8x8 */

	X264DirectPredNone       = 0
	X264DirectPredSpatial    = 1
	X264DirectPredTemporal   = 2
	X264DirectPredAuto       = 3
	X264MeDia                = 0
	X264MeHex                = 1
	X264MeUmh                = 2
	X264MeEsa                = 3
	X264MeTesa               = 4
	X264CqmFlat              = 0
	X264CqmJvt               = 1
	X264CqmCustom            = 2
	X264RcCqp                = 0
	X264RcCrf                = 1
	X264RcAbr                = 2
	X264QpAuto               = 0
	X264AqNone               = 0
	X264AqVariance           = 1
	X264AqAutovariance       = 2
	X264AqAutovarianceBiased = 3
	X264BAdaptNone           = 0
	X264BAdaptFast           = 1
	X264BAdaptTrellis        = 2
	X264WeightpNone          = 0
	X264WeightpSimple        = 1
	X264WeightpSmart         = 2
	X264BPyramidNone         = 0
	X264BPyramidStrict       = 1
	X264BPyramidNormal       = 2
	X264KeyintMinAuto        = 0
	X264KeyintMaxInfinite    = 1 << 30

	/* AVC-Intra flavors */
	X264AvcintraFlavorPanasonic = 0
	X264AvcintraFlavorSony      = 1

	/* !to add missing names */
	/* static const char * const x264_direct_pred_names[] = { "none", "spatial", "temporal", "auto", 0 }; */
	/* static const char * const x264_motion_est_names[] = { "dia", "hex", "umh", "esa", "tesa", 0 }; */
	/* static const char * const x264_b_pyramid_names[] = { "none", "strict", "normal", 0 }; */
	/* static const char * const x264_overscan_names[] = { "undef", "show", "crop", 0 }; */
	/* static const char * const x264_vidformat_names[] = { "component", "pal", "ntsc", "secam", "mac", "undef", 0 }; */
	/* static const char * const x264_fullrange_names[] = { "off", "on", 0 }; */
	/* static const char * const x264_colorprim_names[] = { "", "bt709", "undef", "", "bt470m", "bt470bg", "smpte170m", "smpte240m", "film", "bt2020", "smpte428", "smpte431", "smpte432", 0 }; */
	/* static const char * const x264_transfer_names[] = { "", "bt709", "undef", "", "bt470m", "bt470bg", "smpte170m", "smpte240m", "linear", "log100", "log316", "iec61966-2-4", "bt1361e", "iec61966-2-1", "bt2020-10", "bt2020-12", "smpte2084", "smpte428", "arib-std-b67", 0 }; */
	/* static const char * const x264_colmatrix_names[] = { "GBR", "bt709", "undef", "", "fcc", "bt470bg", "smpte170m", "smpte240m", "YCgCo", "bt2020nc", "bt2020c", "smpte2085", "chroma-derived-nc", "chroma-derived-c", "ICtCp", 0 }; */
	/* static const char * const x264_nal_hrd_names[] = { "none", "vbr", "cbr", 0 }; */
	/* static const char * const x264_avcintra_flavor_names[] = { "panasonic", "sony", 0 }; */

	/* Colorspace type */
	X264CspMask      = 0x00ff /* */
	X264CspNone      = 0x0000 /* Invalid mode     */
	X264CspI400      = 0x0001 /* monochrome 4:0:0 */
	X264CspI420      = 0x0002 /* yuv 4:2:0 planar */
	X264CspYv12      = 0x0003 /* yvu 4:2:0 planar */
	X264CspNv12      = 0x0004 /* yuv 4:2:0, with one y plane and one packed u+v */
	X264CspNv21      = 0x0005 /* yuv 4:2:0, with one y plane and one packed v+u */
	X264CspI422      = 0x0006 /* yuv 4:2:2 planar */
	X264CspYv16      = 0x0007 /* yvu 4:2:2 planar */
	X264CspNv16      = 0x0008 /* yuv 4:2:2, with one y plane and one packed u+v */
	X264CspYuyv      = 0x0009 /* yuyv 4:2:2 packed */
	X264CspUyvy      = 0x000a /* uyvy 4:2:2 packed */
	X264CspV210      = 0x000b /* 10-bit yuv 4:2:2 packed in 32 */
	X264CspI444      = 0x000c /* yuv 4:4:4 planar */
	X264CspYv24      = 0x000d /* yvu 4:4:4 planar */
	X264CspBgr       = 0x000e /* packed bgr 24bits */
	X264CspBgra      = 0x000f /* packed bgr 32bits */
	X264CspRgb       = 0x0010 /* packed rgb 24bits */
	X264CspMax       = 0x0011 /* end of list */
	X264CspVflip     = 0x1000 /* the csp is vertically flipped */
	X264CspHighDepth = 0x2000 /* the csp has a depth of 16 bits per pixel component */

	/* Slice type */
	X264TypeAuto     = 0x0000 /* Let x264 choose the right type */
	X264TypeIdr      = 0x0001
	X264TypeI        = 0x0002
	X264TypeP        = 0x0003
	X264TypeBref     = 0x0004 /* Non-disposable B-frame */
	X264TypeB        = 0x0005
	X264TypeKeyframe = 0x0006 /* IDR or I depending on b_open_gop option */
	/* !to reimplement macro */
	/* #define IS_X264_TYPE_I(x) ((x)==X264_TYPE_I || (x)==X264_TYPE_IDR || (x)==X264_TYPE_KEYFRAME) */
	/* #define IS_X264_TYPE_B(x) ((x)==X264_TYPE_B || (x)==X264_TYPE_BREF) */

	/* Log level */
	X264LogNone    = -1
	X264LogError   = 0
	X264LogWarning = 1
	X264LogInfo    = 2
	X264LogDebug   = 3

	/* Threading */
	X264ThreadsAuto       = 0  /* Automatically select optimal number of threads */
	X264SyncLookaheadAuto = -1 /* Automatically select optimal lookahead thread buffer size */

	/* HRD */
	X264NalHrdNone = 0
	X264NalHrdVbr  = 1
	X264NalHrdCbr  = 2

	ParamBadName  = -1
	ParamBadValue = -2

	// MbinfoConstant.
	MbinfoConstant = 1 << 0
)

// NalUnitType enumeration.
const (
	NalUnknown = int32(iota)
	NalSlice
	NalSliceDpa
	NalSliceDpb
	NalSliceDpc
	NalSliceIdr
	NalSei
	NalSps
	NalPps
	NalAud
	NalFiller
)

// NalPriority enumeration.
const (
	NalPriorityDisposable = int32(iota)
	NalPriorityLow
	NalPriorityHigh
	NalPriorityHighest
)

// PicStruct enumeration.
const (
	PicStructAuto        = int32(iota) // automatically decide (default)
	PicStructProgressive               // progressive frame

	// "TOP" and "BOTTOM" are not supported in x264 (PAFF only)
	PicStructTopBottom       // top field followed by bottom
	PicStructBottomTop       // bottom field followed by top
	PicStructTopBottomTop    // top field, bottom field, top field repeated
	PicStructBottomTopBottom // bottom field, top field, bottom field repeated
	PicStructDouble          // double frame
	PicStructTriple          // triple frame
)

// T opaque handler for encoder.
type T struct{}

// cptr return C pointer.
func (t *T) cptr() *C.x264_t {
	return (*C.x264_t)(unsafe.Pointer(t))
}

// Nal type.
// The data within the payload is already NAL-encapsulated; the ref_idc and type are merely in the struct for easy access by the calling application.
// All data returned in an Nal, including the data in PPayload, is no longer valid after the next call to EncoderEncode.
// Thus it must be used or copied before calling EncoderEncode or EncoderHeaders again.
type Nal struct {
	// NalPriority.
	IRefIdc int32
	// NalUnitType.
	IType int32
	// Start code.
	BLongStartcode int32
	// If this NAL is a slice, the index of the first MB in the slice.
	IFirstMb int32
	// If this NAL is a slice, the index of the last MB in the slice.
	ILastMb int32
	// Size of payload (including any padding) in bytes.
	IPayload int32
	// If param.BAnnexb is set, Annex-B bytestream with startcode.
	// Otherwise, startcode is replaced with a 4-byte size.
	// This size is the size used in mp4/similar muxing; it is equal to IPayload-4.
	PPayload unsafe.Pointer
	// Size of padding in bytes.
	IPadding int32
}

// cptr return C pointer.
func (n *Nal) cptr() *C.x264_nal_t {
	return (*C.x264_nal_t)(unsafe.Pointer(n))
}

// Vui type.
type Vui struct {
	// They will be reduced to be 0 < x <= 65535 and prime.
	ISarHeight int32
	ISarWidth  int32

	// 0=undef, 1=no overscan, 2=overscan.
	IOverscan int32

	// See h264 annex E for the values of the following.
	IVidformat int32
	BFullrange int32
	IColorprim int32
	ITransfer  int32
	IColmatrix int32

	// Both top & bottom.
	IChromaLoc int32
}

// Analyse (encoder analyser parameters) type.
type Analyse struct {
	// Intra partitions.
	Intra uint32
	// Inter partitions.
	Inter uint32

	BTransform8x8 int32
	// Weighting for P-frames.
	IWeightedPred int32
	// Implicit weighting for B-frames.
	BWeightedBipred int32
	// Spatial vs temporal mv prediction.
	IDirectMvPred   int32
	IChromaQpOffset int32

	// Motion estimation algorithm to use (X264_ME_*).
	IMeMethod int32
	// Integer pixel motion estimation search range (from predicted mv).
	IMeRange int32
	// Maximum length of a mv (in pixels). -1 = auto, based on level.
	IMvRange int32
	// Minimum space between threads. -1 = auto, based on number of threads.
	IMvRangeThread int32
	// Subpixel motion estimation quality.
	ISubpelRefine int32
	// Chroma ME for subpel and mode decision in P-frames.
	BChromaMe int32
	// Allow each mb partition to have its own reference number.
	BMixedReferences int32
	// Trellis RD quantization.
	ITrellis int32
	// Early SKIP detection on P-frames.
	BFastPskip int32
	// Transform coefficient thresholding on P-frames.
	BDctDecimate int32
	// Adaptive pseudo-deadzone.
	INoiseReduction int32
	// Psy RD strength.
	FPsyRd float32
	// Psy trellis strength.
	FPsyTrellis float32
	// Toggle all psy optimizations.
	BPsy int32

	// Use input MbInfo data in Picture
	BMbInfo int32
	// Update the values in mb_info according to the results of encoding.
	BMbInfoUpdate int32

	// The deadzone size that will be used in luma quantization {inter, intra}
	ILumaDeadzone [2]int32

	// compute and print PSNR stats
	BPsnr int32
	// compute and print SSIM stats
	BSsim int32
}

// Rc (rate control parameters) type.
type Rc struct {
	// X264_RC_*
	IRcMethod int32

	// 0 to (51 + 6*(x264_bit_depth-8)). 0=lossless.
	IQpConstant int32
	// Min allowed QP value.
	IQpMin int32
	// Max allowed QP value.
	IQpMax int32
	// Max QP step between frames.
	IQpStep int32

	IBitrate int32
	// 1pass VBR, nominal QP.
	FRfConstant float32
	// In CRF mode, maximum CRF as caused by VBV.
	FRfConstantMax float32
	FRateTolerance float32
	IVbvMaxBitrate int32
	IVbvBufferSize int32
	// <=1: fraction of buffer_size. >1: kbit.
	FVbvBufferInit float32
	FIpFactor      float32
	FPbFactor      float32

	// VBV filler: force CBR VBV and use filler bytes to ensure hard-CBR. Implied by NAL-HRD CBR.
	BFiller int32

	// Psy adaptive QP. (X264_AQ_*).
	IAqMode     int32
	FAqStrength float32
	// Macroblock-tree ratecontrol.
	BMbTree    int32
	ILookahead int32

	// 2pass
	// Enable stat writing in psz_stat_out.
	BStatWrite int32
	// Output filename (in UTF-8) of the 2pass stats file.
	PszStatOut *int8
	// Read stat from psz_stat_in and use it.
	BStatRead int32
	_         [4]byte
	// Input filename (in UTF-8) of the 2pass stats file.
	PszStatIn *int8

	// 2pass params (same as ffmpeg ones).
	// 0.0 => cbr, 1.0 => constant qp.
	FQcompress float32
	// Temporally blur quants.
	FQblur float32
	// Temporally blur complexity.
	FComplexityBlur float32
	_               [4]byte
	// Ratecontrol overrides.
	Zones *Zone
	// Number of Zone's.
	IZones int32
	_      [4]byte
	// Alternate method of specifying zones.
	PszZones *int8
}

// CropRect (cropping rectangle parameters) type.
// Added to those implicitly defined by non-mod16 video resolutions.
type CropRect struct {
	Left   uint32
	Top    uint32
	Right  uint32
	Bottom uint32
}

// Zone type.
// Zones: override ratecontrol or other options for specific sections of the video.
// See EncoderReconfig() for which options can be changed.
// If zones overlap, whichever comes later in the list takes precedence.
type Zone struct {
	// Range of frame numbers.
	IStart int32
	// Range of frame numbers.
	IEnd int32
	// Whether to use qp vs bitrate factor.
	BForceQp       int32
	IQp            int32
	FBitrateFactor float32
	_              [4]byte
	Param          *Param
}

// Level (H.264 level restriction information) type.
type Level struct {
	LevelIdc byte
	_        [3]byte
	// Max macroblock processing rate (macroblocks/sec).
	Mbps uint32
	// Max frame size (macroblocks).
	FrameSize uint32
	// Max decoded picture buffer (mbs).
	Dpb uint32
	// Max bitrate (kbit/sec).
	Bitrate uint32
	// Max vbv buffer (kbit).
	Cpb uint32
	// Max vertical mv component range (pixels).
	MvRange uint16
	// Max mvs per 2 consecutive mbs.
	MvsPer2mb byte
	SliceRate byte
	// Min compression ratio.
	Mincr byte
	// Limit bipred to >=8x8.
	Bipred8x8 byte
	// Limit b_direct to >=8x8.
	Direct8x8 byte
	// Forbid interlacing.
	FrameOnly byte
}

// Param type.
type Param struct {
	// CPU flags.
	Cpu uint32
	// Encode multiple frames in parallel.
	IThreads int32
	// Multiple threads for lookahead analysis.
	ILookaheadThreads int32
	// Whether to use slice-based threading.
	BSlicedThreads int32
	// Whether to allow non-deterministic optimizations when threaded.
	BDeterministic int32
	// Force canonical behavior rather than cpu-dependent optimal algorithms.
	BCpuIndependent int32
	// Threaded lookahead buffer.
	ISyncLookahead int32

	// Video Properties.
	IWidth  int32
	IHeight int32
	// CSP of encoded bitstream.
	ICsp      int32
	IBitdepth int32
	ILevelIdc int32
	// Number of frames to encode if known, else 0.
	IFrameTotal int32

	// NAL HRD.
	// Uses Buffering and Picture Timing SEIs to signal HRD. The HRD in H.264 was not designed with VFR in mind.
	// It is therefore not recommendeded to use NAL HRD with VFR.
	// Furthermore, reconfiguring the VBV (via x264_encoder_reconfig) will currently generate invalid HRD.
	INalHrd int32

	Vui Vui

	// Bitstream parameters.
	// Maximum number of reference frames.
	IFrameReference int32
	// Force a DPB size larger than that implied by B-frames and reference frames.
	// Useful in combination with interactive error resilience.
	IDpbSize int32
	// Force an IDR keyframe at this interval.
	IKeyintMax int32
	// Scenecuts closer together than this are coded as I, not IDR.
	IKeyintMin int32
	// How aggressively to insert extra I frames.
	IScenecutThreshold int32
	// Whether or not to use periodic intra refresh instead of IDR frames.
	BIntraRefresh int32

	// How many b-frame between 2 references pictures.
	IBframe         int32
	IBframeAdaptive int32
	IBframeBias     int32
	// Keep some B-frames as references: 0=off, 1=strict hierarchical, 2=normal.
	IBframePyramid int32
	BOpenGop       int32
	BBlurayCompat  int32
	IAvcintraClass int32

	BDeblockingFilter int32
	// [-6, 6] -6 light filter, 6 strong.
	IDeblockingFilterAlphac0 int32
	// [-6, 6]  idem.
	IDeblockingFilterBeta int32

	BCabac        int32
	ICabacInitIdc int32

	BInterlaced       int32
	BConstrainedIntra int32

	ICqmPreset int32
	_          [4]byte
	// Filename (in UTF-8) of CQM file, JM format.
	PszCqmFile *int8

	// Used only if i_cqm_preset == X264CqmCustom.
	Cqm4iy [16]byte
	Cqm4py [16]byte
	Cqm4ic [16]byte
	Cqm4pc [16]byte
	Cqm8iy [64]byte
	Cqm8py [64]byte
	Cqm8ic [64]byte
	Cqm8pc [64]byte

	// Log.
	PfLog       *[0]byte
	PLogPrivate unsafe.Pointer
	ILogLevel   int32
	// Fully reconstruct frames, even when not necessary for encoding. Implied by psz_dump_yuv.
	BFullRecon int32
	// Filename (in UTF-8) for reconstructed frames.
	PszDumpYuv *int8

	// Encoder analyser parameters.
	Analyse Analyse

	_ [4]byte

	// Rate control parameters.
	Rc Rc

	// Cropping Rectangle parameters: added to those implicitly defined by non-mod16 video resolutions.
	CropRect CropRect

	// Frame packing arrangement flag.
	IFramePacking int32

	// Muxing parameters.
	// Generate access unit delimiters.
	BAud int32
	// Put SPS/PPS before each keyframe.
	BRepeatHeaders int32
	// If set, place start codes (4 bytes) before NAL units, otherwise place size (4 bytes) before NAL units.
	BAnnexb int32
	// SPS and PPS id number.
	ISpsId int32
	// VFR input. If 1, use timebase and timestamps for ratecontrol purposes. If 0, use fps only.
	BVfrInput int32
	// Use explicitly set timebase for CFR.
	BPulldown int32
	IFpsNum   uint32
	IFpsDen   uint32
	// Timebase numerator.
	ITimebaseNum uint32
	// Timebase denominator.
	ITimebaseDen uint32

	BTff int32

	// The correct pic_struct must be passed with each input frame.
	// The input timebase should be the timebase corresponding to the output framerate. This should be constant.
	// e.g. for 3:2 pulldown timebase should be 1001/30000.
	// The PTS passed with each frame must be the PTS of the frame after pulldown is applied.
	// Frame doubling and tripling require BVfrInput set to zero (see H.264 Table D-1)
	//
	// Pulldown changes are not clearly defined in H.264. Therefore, it is the calling app's responsibility to manage this.
	BPicStruct int32

	// Used only when b_interlaced=0. Setting this flag makes it possible to flag the stream as PAFF interlaced yet
	// encode all frames progessively. It is useful for encoding 25p and 30p Blu-Ray streams.
	BFakeInterlaced int32

	// Don't optimize header parameters based on video content, e.g. ensure that splitting an input video, compressing
	// each part, and stitching them back together will result in identical SPS/PPS. This is necessary for stitching
	// with container formats that don't allow multiple SPS/PPS.
	BStitchable int32

	// Use OpenCL when available.
	BOpencl int32
	// Specify count of GPU devices to skip, for CLI users.
	IOpenclDevice int32
	_             [4]byte
	// Pass explicit cl_device_id as void*, for API users.
	OpenclDeviceId unsafe.Pointer
	// Filename (in UTF-8) of the compiled OpenCL kernel cache file.
	PszClbinFile *int8

	// Slicing parameters
	// Max size per slice in bytes; includes estimated NAL overhead.
	ISliceMaxSize int32
	// Max number of MBs per slice; overrides i_slice_count.
	ISliceMaxMbs int32
	// Min number of MBs per slice.
	ISliceMinMbs int32
	// Number of slices per frame: forces rectangular slices.
	ISliceCount int32
	// Absolute cap on slices per frame; stops applying slice-max-size and slice-max-mbs if this is reached.
	ISliceCountMax int32

	_           [4]byte
	ParamFree   *[0]byte
	NaluProcess *[0]byte
}

// cptr return C pointer.
func (p *Param) cptr() *C.x264_param_t {
	return (*C.x264_param_t)(unsafe.Pointer(p))
}

// Hrd type.
type Hrd struct {
	CpbInitialArrivalTime float64
	CpbFinalArrivalTime   float64
	CpbRemovalTime        float64
	DpbOutputTime         float64
}

// SeiPayload type.
type SeiPayload struct {
	PayloadSize int32
	PayloadType int32
	Payload     *uint8
}

// Sei type.
type Sei struct {
	NumPayloads int32
	_           [4]byte
	Payloads    *SeiPayload
	SeiFree     *[0]byte
}

// Image type.
type Image struct {
	// Colorspace.
	ICsp int32
	// Number of image planes.
	IPlane int32
	// Strides for each plane.
	IStride [4]int32
	// Pointers to each plane.
	Plane [4]unsafe.Pointer
}

// ImageProperties type.
type ImageProperties struct {
	// In: an array of quantizer offsets to be applied to this image during encoding.
	QuantOffsets *float32
	// In: optional callback to free quant_offsets when used.
	// Useful if one wants to use a different quant_offset array for each frame.
	QuantOffsetsFree *[0]byte

	// In: optional array of flags for each macroblock.
	// Out: if b_mb_info_update is set, x264 will update this array as a result of encoding.
	MbInfo *uint8
	// In: optional callback to free mb_info when used.
	MbInfoFree *[0]byte

	// Out: SSIM of the the frame luma (if x264_param_t.b_ssim is set).
	FSsim float64
	// Out: Average PSNR of the frame (if x264_param_t.b_psnr is set).
	FPsnrAvg float64
	// Out: PSNR of Y, U, and V (if x264_param_t.b_psnr is set).
	FPsnr [3]float64

	// Out: Average effective CRF of the encoded frame.
	FCrfAvg float64
}

// Picture type.
type Picture struct {
	// In: force picture type (if not auto).
	// Out: type of the picture encoded.
	IType int32
	// In: force quantizer for != X264QpAuto.
	IQpplus1 int32
	// In: pic_struct, for pulldown/doubling/etc...used only if b_pic_struct=1.
	// Out: pic_struct element associated with frame.
	IPicStruct int32
	// Out: whether this frame is a keyframe.
	// Important when using modes that result in SEI recovery points being used instead of IDR frames.
	BKeyframe int32
	// In: user pts, Out: pts of encoded picture (user).
	IPts int64
	// Out: frame dts. When the pts of the first frame is close to zero,
	// initial frames may have a negative dts which must be dealt with by any muxer.
	IDts int64
	// In: custom encoding parameters to be set from this frame forwards (in coded order, not display order).
	// If nil, continue using parameters from the previous frame.
	Param *Param
	// In: raw image data.
	// Out: reconstructed image data.
	Img Image
	// In: optional information to modify encoder decisions for this frame.
	// Out: information about the encoded frame.
	Prop ImageProperties
	// Out: HRD timing information. Output only when i_nal_hrd is set.
	Hrdiming Hrd
	// In: arbitrary user SEI (e.g subtitles, AFDs).
	ExtraSei Sei
	// Private user data. copied from input to output frames.
	Opaque *byte
}

// cptr return C pointer.
func (p *Picture) cptr() *C.x264_picture_t {
	return (*C.x264_picture_t)(unsafe.Pointer(p))
}

// NalEncode - encode Nal.
func NalEncode(h *T, dst []byte, nal *Nal) {
	ch := h.cptr()
	cdst := (*C.uint8_t)(unsafe.Pointer(&dst[0]))
	cnal := nal.cptr()

	C.x264_nal_encode(ch, cdst, cnal)
}

// ParamDefault - fill Param with default values and do CPU detection.
func ParamDefault(param *Param) {
	C.x264_param_default(param.cptr())
}

// ParamParse - set one parameter by name. Returns 0 on success.
func ParamParse(param *Param, name string, value string) int32 {
	cparam := param.cptr()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.x264_param_parse(cparam, cname, cvalue)
	v := (int32)(ret)
	return v
}

// ParamDefaultPreset - the same as ParamDefault, but also use the passed preset and tune to modify the default settings
// (either can be nil, which implies no preset or no tune, respectively).
//
// Currently available presets are, ordered from fastest to slowest:
// "ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow", "placebo".
//
// Currently available tunings are:
// "film", "animation", "grain", "stillimage", "psnr", "ssim", "fastdecode", "zerolatency".
//
// Returns 0 on success, negative on failure (e.g. invalid preset/tune name).
func ParamDefaultPreset(param *Param, preset string, tune string) int32 {
	cparam := param.cptr()

	cpreset := C.CString(preset)
	defer C.free(unsafe.Pointer(cpreset))

	ctune := C.CString(tune)
	defer C.free(unsafe.Pointer(ctune))

	ret := C.x264_param_default_preset(cparam, cpreset, ctune)
	v := (int32)(ret)
	return v
}

// ParamApplyFastfirstpass - if first-pass mode is set (rc.b_stat_read == 0, rc.b_stat_write == 1),
// modify the encoder settings to disable options generally not useful on the first pass.
func ParamApplyFastfirstpass(param *Param) {
	cparam := param.cptr()
	C.x264_param_apply_fastfirstpass(cparam)
}

// ParamApplyProfile - applies the restrictions of the given profile.
//
// Currently available profiles are, from most to least restrictive:
// "baseline", "main", "high", "high10", "high422", "high444".
// (can be nil, in which case the function will do nothing).
//
// Returns 0 on success, negative on failure (e.g. invalid profile name).
func ParamApplyProfile(param *Param, profile string) int32 {
	cparam := param.cptr()

	cprofile := C.CString(profile)
	defer C.free(unsafe.Pointer(cprofile))

	ret := C.x264_param_apply_profile(cparam, cprofile)
	v := (int32)(ret)
	return v
}

// PictureInit - initialize an Picture. Needs to be done if the calling application
// allocates its own Picture as opposed to using PictureAlloc.
func PictureInit(pic *Picture) {
	cpic := pic.cptr()
	C.x264_picture_init(cpic)
}

// PictureAlloc - alloc data for a Picture. You must call PictureClean on it.
// Returns 0 on success, or -1 on malloc failure or invalid colorspace.
func PictureAlloc(pic *Picture, iCsp int32, iWidth int32, iHeight int32) int32 {
	cpic := pic.cptr()

	ciCsp := (C.int)(iCsp)
	ciWidth := (C.int)(iWidth)
	ciHeight := (C.int)(iHeight)

	ret := C.x264_picture_alloc(cpic, ciCsp, ciWidth, ciHeight)
	v := (int32)(ret)
	return v
}

// PictureClean - free associated resource for a Picture allocated with PictureAlloc ONLY.
func PictureClean(pic *Picture) {
	cpic := pic.cptr()
	C.x264_picture_clean(cpic)
}

// EncoderOpen - create a new encoder handler, all parameters from Param are copied.
func EncoderOpen(param *Param) *T {
	cparam := param.cptr()

	ret := C.x264_encoder_open(cparam)
	v := *(**T)(unsafe.Pointer(&ret))
	return v
}

// EncoderReconfig - various parameters from Param are copied.
// Returns 0 on success, negative on parameter validation error.
func EncoderReconfig(enc *T, param *Param) int32 {
	cenc := enc.cptr()
	cparam := param.cptr()

	ret := C.x264_encoder_reconfig(cenc, cparam)
	v := (int32)(ret)
	return v
}

// EncoderParameters - copies the current internal set of parameters to the pointer provided.
func EncoderParameters(enc *T, param *Param) {
	cenc := enc.cptr()
	cparam := param.cptr()

	C.x264_encoder_parameters(cenc, cparam)
}

// EncoderHeaders - return the SPS and PPS that will be used for the whole stream.
// Returns the number of bytes in the returned NALs or negative on error.
func EncoderHeaders(enc *T, ppNal []*Nal, piNal *int32) int32 {
	cenc := enc.cptr()

	cppNal := (**C.x264_nal_t)(unsafe.Pointer(&ppNal[0]))
	cpiNal := (*C.int)(unsafe.Pointer(piNal))

	ret := C.x264_encoder_headers(cenc, cppNal, cpiNal)
	v := (int32)(ret)
	return v
}

// EncoderEncode - encode one picture.
// Returns the number of bytes in the returned NALs, negative on error and zero if no NAL units returned.
func EncoderEncode(enc *T, ppNal []*Nal, piNal *int32, picIn *Picture, picOut *Picture) int32 {
	cenc := enc.cptr()

	cppNal := (**C.x264_nal_t)(unsafe.Pointer(&ppNal[0]))
	cpiNal := (*C.int)(unsafe.Pointer(piNal))

	cpicIn := picIn.cptr()
	cpicOut := picOut.cptr()

	ret := C.x264_encoder_encode(cenc, cppNal, cpiNal, cpicIn, cpicOut)
	v := (int32)(ret)
	return v
}

// EncoderClose - close an encoder handler.
func EncoderClose(enc *T) {
	cenc := enc.cptr()
	C.x264_encoder_close(cenc)
}

// EncoderDelayedFrames - return the number of currently delayed (buffered) frames.
// This should be used at the end of the stream, to know when you have all the encoded frames.
func EncoderDelayedFrames(enc *T) int32 {
	cenc := enc.cptr()

	ret := C.x264_encoder_delayed_frames(cenc)
	v := (int32)(ret)
	return v
}

// EncoderMaximumDelayedFrames - return the maximum number of delayed (buffered) frames that can occur with the current parameters.
func EncoderMaximumDelayedFrames(enc *T) int32 {
	cenc := enc.cptr()

	ret := C.x264_encoder_maximum_delayed_frames(cenc)
	v := (int32)(ret)
	return v
}

// EncoderIntraRefresh - If an intra refresh is not in progress, begin one with the next P-frame.
// If an intra refresh is in progress, begin one as soon as the current one finishes.
// Requires that BIntraRefresh be set.
//
// Should not be called during an x264_encoder_encode.
func EncoderIntraRefresh(enc *T) {
	cenc := enc.cptr()
	C.x264_encoder_intra_refresh(cenc)
}

// EncoderInvalidateReference - An interactive error resilience tool, designed for use in a low-latency one-encoder-few-clients system.
// Should not be called during an EncoderEncode, but multiple calls can be made simultaneously.
//
// Returns 0 on success, negative on failure.
func EncoderInvalidateReference(enc *T, pts int) int32 {
	cenc := enc.cptr()
	cpts := (C.int64_t)(pts)

	ret := C.x264_encoder_invalidate_reference(cenc, cpts)
	v := (int32)(ret)
	return v
}
