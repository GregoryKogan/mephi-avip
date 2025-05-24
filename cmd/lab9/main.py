import numpy as np
import librosa
import librosa.display
import scipy.signal as signal
import matplotlib.pyplot as plt
from scipy.ndimage import uniform_filter1d
import soundfile as sf

# Parameters
N_FFT = 2048
HOP_LENGTH = 512
WINDOW = "hann"
FRAME_STEP = 0.1  # seconds
FREQ_STEP = 50  # Hz


def load_audio(path):
    y, sr = librosa.load(path, sr=None, mono=True)
    return y, sr


def compute_spectrogram(y, sr):
    S = np.abs(librosa.stft(y, n_fft=N_FFT, hop_length=HOP_LENGTH, window=WINDOW))
    freqs = librosa.fft_frequencies(sr=sr, n_fft=N_FFT)
    times = librosa.frames_to_time(np.arange(S.shape[1]), sr=sr, hop_length=HOP_LENGTH)
    return S, freqs, times


def save_spectrogram(S, freqs, times, out_path, title):
    plt.figure(figsize=(10, 6))
    librosa.display.specshow(
        librosa.amplitude_to_db(S, ref=np.max),
        x_coords=times,
        y_coords=freqs,
        x_axis="time",
        y_axis="log",
        cmap="magma",
    )
    plt.colorbar(format="%+2.0f dB")
    plt.title(title)
    plt.tight_layout()
    plt.savefig(out_path)
    plt.close()


def reduce_noise(S, method="wiener"):
    if method == "wiener":
        return signal.wiener(S)
    elif method == "savgol":
        return signal.savgol_filter(S, window_length=11, polyorder=2, axis=1)
    elif method == "lowpass":
        return uniform_filter1d(S, size=5, axis=1)
    else:
        return S


def find_energy_peaks(S, freqs, times, sr):
    frame_size = int(FRAME_STEP * sr / HOP_LENGTH)
    freq_bin = int(FREQ_STEP / (freqs[1] - freqs[0]))
    energy_peaks = []
    for i in range(0, S.shape[1] - frame_size, frame_size):
        patch = S[:, i : i + frame_size]
        E = patch.sum(axis=1)
        max_idx = np.argmax(E)
        energy_peaks.append((times[i], freqs[max_idx], E[max_idx]))
    return energy_peaks


def plot_energy_peaks(peaks, out_path):
    times = [p[0] for p in peaks]
    freqs = [p[1] for p in peaks]
    energies = [p[2] for p in peaks]
    plt.figure(figsize=(10, 4))
    sc = plt.scatter(times, freqs, c=energies, cmap="hot", s=50)
    plt.colorbar(sc, label="Energy")
    plt.xlabel("Time (s)")
    plt.ylabel("Frequency (Hz)")
    plt.title("Energy Peaks Over Time")
    plt.tight_layout()
    plt.savefig(out_path)
    plt.close()


def restore_audio(S, phase):
    S_complex = S * np.exp(1j * phase)
    y = librosa.istft(S_complex, hop_length=HOP_LENGTH, window=WINDOW)
    return y


if __name__ == "__main__":
    path = "assets/sample.wav"
    y, sr = load_audio(path)
    S, freqs, times = compute_spectrogram(y, sr)
    save_spectrogram(
        S, freqs, times, "output/lab9/original_spectrogram.png", "Original Spectrogram"
    )

    for method in ["wiener", "savgol", "lowpass"]:
        S_denoised = reduce_noise(S, method=method)
        save_spectrogram(
            S_denoised,
            freqs,
            times,
            f"output/lab9/denoised_{method}_spectrogram.png",
            "Denoised Spectrogram",
        )

    S_denoised = reduce_noise(S, method="wiener")
    phase = np.angle(librosa.stft(y, n_fft=N_FFT, hop_length=HOP_LENGTH, window=WINDOW))
    y_restored = restore_audio(S_denoised, phase)
    sf.write("output/lab9/restored.wav", y_restored, sr)

    peaks = find_energy_peaks(S, freqs, times, sr)
    plot_energy_peaks(peaks, "output/lab9/energy_peaks.png")
